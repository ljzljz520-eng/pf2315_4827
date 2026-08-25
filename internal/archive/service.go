package archive

import (
	"fmt"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Service struct {
	store       *store.Store
	auditPrefix string
	nextAudit   int
}

func New(s *store.Store, prefix string) *Service {
	return &Service{store: s, auditPrefix: prefix, nextAudit: 1}
}

func (s *Service) Archive(recordID, actor string, now time.Time) (model.Record, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if err := record.ValidateForTransition(model.StatusArchived); err != nil {
		return model.Record{}, err
	}
	previous := record.Status
	record.Status = model.StatusArchived
	record.ArchivedAt = now
	record.UpdatedAt = now
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	event := model.AuditEvent{ID: fmt.Sprintf("%s-archive-%03d", s.auditPrefix, s.nextAudit), RecordID: recordID, Action: "archive", Actor: actor, From: previous, To: model.StatusArchived, Detail: "归档", CreatedAt: now}
	s.nextAudit++
	if err := s.store.PutAudit(event); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) ArchiveEdition(edition, actor string, now time.Time) (model.ArchiveBatch, error) {
	records, err := s.store.ListRecords(model.SearchFilter{Edition: edition})
	if err != nil {
		return model.ArchiveBatch{}, err
	}
	batch := model.ArchiveBatch{Edition: edition, Actor: actor, ArchivedAt: now, RecordIDs: make([]string, 0)}
	for _, record := range records {
		if record.Status == model.StatusApproved {
			if _, err := s.Archive(record.ID, actor, now); err != nil {
				return model.ArchiveBatch{}, err
			}
			batch.RecordIDs = append(batch.RecordIDs, record.ID)
		}
	}
	return batch, nil
}

func (s *Service) RestoreRejected(recordID, actor string, now time.Time) (model.Record, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if err := record.ValidateForTransition(model.StatusDraft); err != nil {
		return model.Record{}, err
	}
	previous := record.Status
	record.Status = model.StatusDraft
	record.UpdatedAt = now
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	event := model.AuditEvent{ID: fmt.Sprintf("%s-restore-%03d", s.auditPrefix, s.nextAudit), RecordID: recordID, Action: "restore", Actor: actor, From: previous, To: model.StatusDraft, Detail: "退回修改", CreatedAt: now}
	s.nextAudit++
	if err := s.store.PutAudit(event); err != nil {
		return model.Record{}, err
	}
	return record, nil
}
