package review

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Service struct {
	store *store.Store
	ids   *IDGenerator
}

type IDGenerator struct {
	prefix string
	next   int
}

func New(s *store.Store, prefix string) *Service {
	return &Service{store: s, ids: &IDGenerator{prefix: prefix, next: 1}}
}

func (g *IDGenerator) Next() string {
	id := fmt.Sprintf("%s-audit-%04d", g.prefix, g.next)
	g.next++
	return id
}

func (s *Service) Submit(recordID, actor string, now time.Time) (model.Record, error) {
	return s.transition(recordID, model.StatusSubmitted, actor, "提交审核", now)
}

func (s *Service) Decide(decision model.ReviewDecision, now time.Time) (model.Record, error) {
	if strings.TrimSpace(decision.Reviewer) == "" {
		return model.Record{}, fmt.Errorf("reviewer is required")
	}
	next := model.StatusRejected
	detail := strings.TrimSpace(decision.Reason)
	if decision.Approved {
		next, detail = model.StatusApproved, "审核通过"
	}
	if detail == "" {
		detail = "审核退回"
	}
	return s.transition(decision.RecordID, next, decision.Reviewer, detail, now)
}

func (s *Service) transition(recordID string, next model.Status, actor, detail string, now time.Time) (model.Record, error) {
	record, err := s.store.GetRecord(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if err := record.ValidateForTransition(next); err != nil {
		return model.Record{}, err
	}
	previous := record.Status
	record.Status = next
	record.UpdatedAt = now
	if next == model.StatusApproved {
		record.PublishedAt = now
	}
	if next == model.StatusRejected {
		record.Notes = append(record.Notes, detail)
	}
	if err := s.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	event := model.AuditEvent{ID: s.ids.Next(), RecordID: recordID, Action: "status_transition", Actor: actor, From: previous, To: next, Detail: detail, CreatedAt: now}
	if err := s.store.PutAudit(event); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) History(recordID string) ([]model.AuditEvent, error) {
	return s.store.ListAudits(recordID)
}

func (s *Service) CanReview(record model.Record) bool { return record.Status == model.StatusSubmitted }
