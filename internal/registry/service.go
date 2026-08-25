package registry

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Service struct {
	store  *store.Store
	prefix string
	next   int
}

func New(s *store.Store, prefix string) *Service { return &Service{store: s, prefix: prefix, next: 1} }

func (s *Service) Store() *store.Store { return s.store }

func (s *Service) nextID() string {
	id := fmt.Sprintf("%s-%04d", s.prefix, s.next)
	s.next++
	return id
}

func (s *Service) Register(title, edition, summary, ageRange string, materials, steps []string, editor string, now time.Time) (model.Record, error) {
	record := model.NewRecord(s.nextID(), title, edition, summary, ageRange, materials, steps, now)
	record.Editor = strings.TrimSpace(editor)
	record = model.NormalizeRecord(record)
	if issues := record.Validate(); len(issues) > 0 {
		return model.Record{}, fmt.Errorf("record validation failed: %s", formatIssues(issues))
	}
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Load(id string) (model.Record, error) { return s.store.GetRecord(id) }

func (s *Service) Replace(record model.Record) error {
	record = model.NormalizeRecord(record)
	if issues := record.Validate(); len(issues) > 0 {
		return fmt.Errorf("record validation failed: %s", formatIssues(issues))
	}
	return s.store.PutRecord(record)
}

func (s *Service) Amend(id, editor, title, summary string, materials, steps []string, now time.Time) (model.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusDraft && record.Status != model.StatusRejected {
		return model.Record{}, fmt.Errorf("record %s cannot be amended in %s", id, record.Status)
	}
	if title != "" {
		record.Title = strings.TrimSpace(title)
	}
	if summary != "" {
		record.Summary = strings.TrimSpace(summary)
	}
	if materials != nil {
		record.Materials = append([]string(nil), materials...)
	}
	if steps != nil {
		record.Steps = append([]string(nil), steps...)
	}
	record.Editor = strings.TrimSpace(editor)
	record.Version++
	record.UpdatedAt = now
	if issues := record.Validate(); len(issues) > 0 {
		return model.Record{}, fmt.Errorf("record validation failed: %s", formatIssues(issues))
	}
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) AddNote(id, note string, now time.Time) (model.Record, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	note = strings.TrimSpace(note)
	if note == "" {
		return model.Record{}, fmt.Errorf("note is required")
	}
	record.Notes = append(record.Notes, note)
	record.UpdatedAt = now
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) ValidateForSubmission(id string) ([]model.ValidationIssue, error) {
	record, err := s.store.GetRecord(id)
	if err != nil {
		return nil, err
	}
	return record.Validate(), nil
}

func formatIssues(issues []model.ValidationIssue) string {
	parts := make([]string, len(issues))
	for i, issue := range issues {
		parts[i] = issue.Field + ": " + issue.Message
	}
	return strings.Join(parts, ", ")
}
