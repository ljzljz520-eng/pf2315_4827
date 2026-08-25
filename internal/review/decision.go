package review

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
)

type DecisionSummary struct {
	RecordID string
	Previous model.Status
	Current  model.Status
	Reviewer string
	Reason   string
	At       time.Time
}

func ValidateDecision(decision model.ReviewDecision) error {
	if strings.TrimSpace(decision.RecordID) == "" {
		return fmt.Errorf("record id is required")
	}
	if strings.TrimSpace(decision.Reviewer) == "" {
		return fmt.Errorf("reviewer is required")
	}
	if !decision.Approved && strings.TrimSpace(decision.Reason) == "" {
		return fmt.Errorf("rejection reason is required")
	}
	return nil
}

func DecisionLabel(decision model.ReviewDecision) string {
	if decision.Approved {
		return "approved"
	}
	return "rejected"
}

func (s *Service) DecideWithSummary(decision model.ReviewDecision, now time.Time) (model.Record, DecisionSummary, error) {
	if err := ValidateDecision(decision); err != nil {
		return model.Record{}, DecisionSummary{}, err
	}
	before, err := s.store.GetRecord(decision.RecordID)
	if err != nil {
		return model.Record{}, DecisionSummary{}, err
	}
	after, err := s.Decide(decision, now)
	if err != nil {
		return model.Record{}, DecisionSummary{}, err
	}
	return after, DecisionSummary{RecordID: after.ID, Previous: before.Status, Current: after.Status, Reviewer: decision.Reviewer, Reason: decision.Reason, At: now}, nil
}

func (s *Service) NeedsSecondReview(record model.Record) bool {
	return record.Status == model.StatusRejected && len(record.Notes) > 0
}

func (s *Service) ReasonFor(record model.Record) string {
	if len(record.Notes) == 0 {
		return ""
	}
	return strings.TrimSpace(record.Notes[len(record.Notes)-1])
}
