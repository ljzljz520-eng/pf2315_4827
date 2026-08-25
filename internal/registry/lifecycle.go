package registry

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
)

type LifecycleEvent struct {
	Status  model.Status
	Content string
	Actor   string
	At      time.Time
}

func (s *Service) ReplaceStatusContent(id, actor, content string, now time.Time) (model.Record, error) {
	record, err := s.Load(id)
	if err != nil {
		return model.Record{}, err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return model.Record{}, fmt.Errorf("status content is required")
	}
	if len(record.Notes) == 0 {
		record.Notes = append(record.Notes, content)
	} else {
		record.Notes[len(record.Notes)-1] = content
	}
	record.UpdatedAt = now
	if actor != "" {
		record.Editor = strings.TrimSpace(actor)
	}
	if err := s.Replace(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) RecordStatusSnapshot(id string) (string, error) {
	record, err := s.Load(id)
	if err != nil {
		return "", err
	}
	if len(record.Notes) == 0 {
		return "", nil
	}
	return record.Notes[len(record.Notes)-1], nil
}

func (s *Service) EligibleForEditorialUpdate(record model.Record) bool {
	return record.Status == model.StatusDraft || record.Status == model.StatusRejected
}

func (s *Service) EligibleForPublicListing(record model.Record) bool {
	return record.Status == model.StatusApproved || record.Status == model.StatusArchived
}

func (s *Service) VersionLabel(record model.Record) string {
	return fmt.Sprintf("%s-v%d", record.ID, record.Version)
}
