package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"example.com/scienceweekly/internal/model"
	"go.etcd.io/bbolt"
)

type ExportBundle struct {
	Records     []model.Record
	Audits      []model.AuditEvent
	Workflows   []model.Workflow
	Attachments []model.Attachment
}

func (s *Store) ExportBundle(filter model.SearchFilter) (ExportBundle, error) {
	records, err := s.ListRecords(filter)
	if err != nil {
		return ExportBundle{}, err
	}
	bundle := ExportBundle{Records: records, Audits: make([]model.AuditEvent, 0), Workflows: make([]model.Workflow, 0), Attachments: make([]model.Attachment, 0)}
	for _, record := range records {
		audits, auditErr := s.ListAudits(record.ID)
		if auditErr != nil {
			return ExportBundle{}, auditErr
		}
		bundle.Audits = append(bundle.Audits, audits...)
		workflows, workflowErr := s.ListWorkflows(record.ID)
		if workflowErr != nil {
			return ExportBundle{}, workflowErr
		}
		bundle.Workflows = append(bundle.Workflows, workflows...)
		attachments, attachmentErr := s.ListAttachments(record.ID)
		if attachmentErr != nil {
			return ExportBundle{}, attachmentErr
		}
		bundle.Attachments = append(bundle.Attachments, attachments...)
	}
	return bundle, nil
}

func EncodeBundle(bundle ExportBundle) ([]byte, error) {
	return json.MarshalIndent(bundle, "", "  ")
}

func (s *Store) ImportBundle(bundle ExportBundle) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, record := range bundle.Records {
			data, err := marshal(record)
			if err != nil {
				return err
			}
			if err := tx.Bucket(recordBucket).Put(key(record.ID), data); err != nil {
				return err
			}
		}
		for _, event := range bundle.Audits {
			data, err := marshal(event)
			if err != nil {
				return err
			}
			if err := tx.Bucket(auditBucket).Put(key(event.ID), data); err != nil {
				return err
			}
		}
		for _, workflow := range bundle.Workflows {
			data, err := marshal(workflow)
			if err != nil {
				return err
			}
			if err := tx.Bucket(workflowBucket).Put(key(workflow.ID), data); err != nil {
				return err
			}
		}
		for _, attachment := range bundle.Attachments {
			data, err := marshal(attachment)
			if err != nil {
				return err
			}
			if err := tx.Bucket(attachmentBucket).Put(key(attachment.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ExportLines(filter model.SearchFilter) ([]string, error) {
	records, err := s.ListRecords(filter)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(records))
	for _, record := range records {
		lines = append(lines, strings.Join([]string{record.ID, record.Edition, record.Title, record.Summary, string(record.Status)}, "\t"))
	}
	sort.Strings(lines)
	return lines, nil
}

func DecodeBundle(data []byte) (ExportBundle, error) {
	var bundle ExportBundle
	if len(bytes.TrimSpace(data)) == 0 {
		return bundle, fmt.Errorf("bundle is empty")
	}
	if err := json.Unmarshal(data, &bundle); err != nil {
		return bundle, fmt.Errorf("decode bundle: %w", err)
	}
	return bundle, nil
}
