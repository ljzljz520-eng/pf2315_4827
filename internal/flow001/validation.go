package flow001

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/review"
)

type BatchCheck struct {
	RecordID string
	Valid    bool
	Issues   []model.ValidationIssue
}

func (s *Service) ValidateRecord(recordID string) (BatchCheck, error) {
	record, err := s.registry.Load(recordID)
	if err != nil {
		return BatchCheck{}, err
	}
	issues := record.Validate()
	if policyIssues := review.DefaultPolicy().Check(record); len(policyIssues) > 0 {
		issues = append(issues, policyIssues...)
	}
	return BatchCheck{RecordID: record.ID, Valid: len(issues) == 0, Issues: issues}, nil
}

func (s *Service) ValidateRecords(recordIDs []string) []BatchCheck {
	checks := make([]BatchCheck, 0, len(recordIDs))
	for _, recordID := range recordIDs {
		check, err := s.ValidateRecord(recordID)
		if err != nil {
			checks = append(checks, BatchCheck{RecordID: recordID, Valid: false, Issues: []model.ValidationIssue{{Field: "record", Message: err.Error()}}})
			continue
		}
		checks = append(checks, check)
	}
	return checks
}

func (s *Service) RequireValidRecord(recordID string) error {
	check, err := s.ValidateRecord(recordID)
	if err != nil {
		return err
	}
	if check.Valid {
		return nil
	}
	parts := make([]string, len(check.Issues))
	for i, issue := range check.Issues {
		parts[i] = issue.Field + ":" + issue.Message
	}
	return fmt.Errorf("record %s invalid: %s", recordID, strings.Join(parts, "; "))
}

func (s *Service) SetFixedTime(now time.Time) { s.now = now }

func (s *Service) StatusEntries(recordID string) ([]string, error) {
	record, err := s.registry.Load(recordID)
	if err != nil {
		return nil, err
	}
	entries := make([]string, len(record.Notes))
	copy(entries, record.Notes)
	return entries, nil
}

func (s *Service) LatestStatus(recordID string) (string, error) {
	entries, err := s.StatusEntries(recordID)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", nil
	}
	return entries[len(entries)-1], nil
}

func (s *Service) StatusAt(recordID string, index int) (string, error) {
	entries, err := s.StatusEntries(recordID)
	if err != nil {
		return "", err
	}
	if index < 0 || index >= len(entries) {
		return "", fmt.Errorf("status index %d is outside range", index)
	}
	return entries[index], nil
}

func (s *Service) CompareStatusEntries(recordID string, expected []string) (bool, error) {
	actual, err := s.StatusEntries(recordID)
	if err != nil {
		return false, err
	}
	if len(actual) != len(expected) {
		return false, nil
	}
	for index := range actual {
		if actual[index] != expected[index] {
			return false, nil
		}
	}
	return true, nil
}
