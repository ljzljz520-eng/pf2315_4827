package flow001

import (
	"fmt"
	"time"

	"example.com/scienceweekly/internal/archive"
	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/review"
	"example.com/scienceweekly/internal/store"
)

type Service struct {
	registry *registry.Service
	review   *review.Service
	archive  *archive.Service
	store    *store.Store
	now      time.Time
}

func New(r *registry.Service, v *review.Service, a *archive.Service, now time.Time) *Service {
	return &Service{registry: r, review: v, archive: a, store: r.Store(), now: now}
}

func (s *Service) CreateReviewArchive(title, edition, summary, age string, materials, steps []string, editor, reviewer string) (model.Record, error) {
	record, err := s.registry.Register(title, edition, summary, age, materials, steps, editor, s.now)
	if err != nil {
		return model.Record{}, err
	}
	if err := s.reviewPolicy(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.review.Submit(record.ID, editor, s.now); err != nil {
		return model.Record{}, err
	}
	record, err = s.review.Decide(model.ReviewDecision{RecordID: record.ID, Approved: true, Reviewer: reviewer}, s.now)
	if err != nil {
		return model.Record{}, err
	}
	return s.archive.Archive(record.ID, reviewer, s.now)
}

func (s *Service) reviewPolicy(record model.Record) error {
	policy := review.DefaultPolicy()
	if !policy.Allows(record) {
		return fmt.Errorf("review policy rejected %s: %s", record.ID, policy.Explain(record))
	}
	return nil
}

func (s *Service) SubmitAndDecide(recordID, reviewer string, approved bool, reason string) (model.Record, error) {
	if _, err := s.review.Submit(recordID, reviewer, s.now); err != nil {
		return model.Record{}, err
	}
	return s.review.Decide(model.ReviewDecision{RecordID: recordID, Approved: approved, Reviewer: reviewer, Reason: reason}, s.now)
}

func (s *Service) ArchiveApproved(recordID, actor string) (model.Record, error) {
	return s.archive.Archive(recordID, actor, s.now)
}

func (s *Service) AddStatus(recordID, statusContent string) (model.Record, error) {
	record, err := s.registry.Load(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if statusContent == "" {
		return model.Record{}, fmt.Errorf("status content is required")
	}
	if len(record.Notes) > 0 {
		statusContent = record.Notes[len(record.Notes)-1]
	}
	record.Notes = append(record.Notes, statusContent)
	record.UpdatedAt = s.now
	if err := s.registry.Replace(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) StatusHistory(recordID string) (model.Record, error) {
	return s.registry.Load(recordID)
}
