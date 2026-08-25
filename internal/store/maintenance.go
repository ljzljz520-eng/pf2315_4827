package store

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"go.etcd.io/bbolt"
)

type IntegrityReport struct {
	Records     int
	Audits      int
	Workflows   int
	Attachments int
	Errors      []string
}

func (s *Store) CheckIntegrity() (IntegrityReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return IntegrityReport{}, err
	}
	report := IntegrityReport{Errors: make([]string, 0)}
	err := s.db.View(func(tx *bbolt.Tx) error {
		records := tx.Bucket(recordBucket)
		_ = records.ForEach(func(k, value []byte) error {
			var record model.Record
			if err := unmarshal(value, &record); err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("record %s: %v", k, err))
				return nil
			}
			report.Records++
			if record.ID != string(k) {
				report.Errors = append(report.Errors, fmt.Sprintf("record key mismatch %s", k))
			}
			return nil
		})
		audits := tx.Bucket(auditBucket)
		_ = audits.ForEach(func(k, value []byte) error {
			var event model.AuditEvent
			if err := unmarshal(value, &event); err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("audit %s: %v", k, err))
				return nil
			}
			report.Audits++
			return nil
		})
		workflows := tx.Bucket(workflowBucket)
		_ = workflows.ForEach(func(k, value []byte) error {
			var workflow model.Workflow
			if err := unmarshal(value, &workflow); err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("workflow %s: %v", k, err))
				return nil
			}
			report.Workflows++
			if err := workflow.Validate(); err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("workflow %s: %v", k, err))
			}
			return nil
		})
		attachments := tx.Bucket(attachmentBucket)
		_ = attachments.ForEach(func(k, value []byte) error {
			var attachment model.Attachment
			if err := unmarshal(value, &attachment); err != nil {
				report.Errors = append(report.Errors, fmt.Sprintf("attachment %s: %v", k, err))
				return nil
			}
			report.Attachments++
			return nil
		})
		return nil
	})
	return report, err
}

func (s *Store) SaveWorkflow(workflow model.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}
	return s.PutWorkflow(workflow)
}

func (s *Store) AdvanceWorkflow(id string, now time.Time) (model.Workflow, error) {
	workflow, err := s.GetWorkflow(id)
	if err != nil {
		return model.Workflow{}, err
	}
	if err := workflow.Advance(now); err != nil {
		return model.Workflow{}, err
	}
	if err := s.PutWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (s *Store) Attach(recordID, name, mediaType string, content []byte, checksum string, now time.Time) (model.Attachment, error) {
	if strings.TrimSpace(recordID) == "" {
		return model.Attachment{}, fmt.Errorf("record id is required")
	}
	attachment := model.Attachment{ID: fmt.Sprintf("%s-attachment-%d", recordID, len(content)), RecordID: recordID, Name: strings.TrimSpace(name), MediaType: strings.TrimSpace(mediaType), Content: append([]byte(nil), content...), Checksum: checksum, CreatedAt: now}
	if err := attachment.Validate(); err != nil {
		return model.Attachment{}, err
	}
	if err := s.PutAttachment(attachment); err != nil {
		return model.Attachment{}, err
	}
	return attachment, nil
}

func (s *Store) DeleteArchivedBefore(cutoff time.Time) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return 0, err
	}
	removed := 0
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		keys := sortedKeys(bucket)
		for _, k := range keys {
			var record model.Record
			if err := unmarshal(bucket.Get(k), &record); err != nil {
				return err
			}
			if record.Status == model.StatusArchived && !record.ArchivedAt.After(cutoff) {
				if err := bucket.Delete(k); err != nil {
					return err
				}
				removed++
			}
		}
		return nil
	})
	return removed, err
}

func (s *Store) PurgeRecordArtifacts(recordID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucketName := range [][]byte{auditBucket, workflowBucket, attachmentBucket} {
			bucket := tx.Bucket(bucketName)
			keys := sortedKeys(bucket)
			for _, key := range keys {
				value := bucket.Get(key)
				if bytesContainRecord(value, recordID) {
					if err := bucket.Delete(key); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func bytesContainRecord(value []byte, recordID string) bool {
	return strings.Contains(string(value), fmt.Sprintf(`"record_id":"%s"`, recordID))
}
