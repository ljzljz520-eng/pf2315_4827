package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusSubmitted Status = "submitted"
	StatusApproved  Status = "approved"
	StatusArchived  Status = "archived"
	StatusRejected  Status = "rejected"
)

type Record struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Edition     string    `json:"edition"`
	Summary     string    `json:"summary"`
	AgeRange    string    `json:"age_range"`
	Materials   []string  `json:"materials"`
	Steps       []string  `json:"steps"`
	Status      Status    `json:"status"`
	Version     int       `json:"version"`
	Editor      string    `json:"editor"`
	Reviewer    string    `json:"reviewer"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
	ArchivedAt  time.Time `json:"archived_at"`
	Notes       []string  `json:"notes"`
}

type AuditEvent struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	From      Status    `json:"from"`
	To        Status    `json:"to"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type Workflow struct {
	ID          string    `json:"id"`
	RecordID    string    `json:"record_id"`
	Name        string    `json:"name"`
	Owner       string    `json:"owner"`
	State       string    `json:"state"`
	CurrentStep int       `json:"current_step"`
	Steps       []string  `json:"steps"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
}

type Attachment struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Name      string    `json:"name"`
	MediaType string    `json:"media_type"`
	Content   []byte    `json:"content"`
	Checksum  string    `json:"checksum"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchFilter struct {
	Text            string
	Status          Status
	Edition         string
	AgeRange        string
	IncludeArchived bool
}

type ImportRow struct {
	Title     string
	Edition   string
	Summary   string
	AgeRange  string
	Materials string
	Steps     string
}

type ValidationIssue struct {
	Field   string
	Message string
}

func (r Record) Validate() []ValidationIssue {
	issues := make([]ValidationIssue, 0, 5)
	if strings.TrimSpace(r.ID) == "" {
		issues = append(issues, ValidationIssue{"id", "id is required"})
	}
	if strings.TrimSpace(r.Title) == "" {
		issues = append(issues, ValidationIssue{"title", "title is required"})
	}
	if strings.TrimSpace(r.Edition) == "" {
		issues = append(issues, ValidationIssue{"edition", "edition is required"})
	}
	if strings.TrimSpace(r.Summary) == "" {
		issues = append(issues, ValidationIssue{"summary", "summary is required"})
	}
	if len(r.Materials) < 1 {
		issues = append(issues, ValidationIssue{"materials", "at least one material is required"})
	}
	if len(r.Steps) < 2 {
		issues = append(issues, ValidationIssue{"steps", "at least two steps are required"})
	}
	if r.Status == "" {
		issues = append(issues, ValidationIssue{"status", "status is required"})
	}
	return issues
}

func (r Record) ValidateForTransition(next Status) error {
	if next == r.Status {
		return errors.New("status is unchanged")
	}
	allowed := map[Status][]Status{
		StatusDraft:     {StatusSubmitted},
		StatusSubmitted: {StatusApproved, StatusRejected, StatusDraft},
		StatusApproved:  {StatusArchived},
		StatusRejected:  {StatusDraft},
		StatusArchived:  {},
	}
	for _, candidate := range allowed[r.Status] {
		if candidate == next {
			return nil
		}
	}
	return fmt.Errorf("cannot transition from %s to %s", r.Status, next)
}

func (r Record) Clone() Record {
	copy := r
	copy.Materials = append([]string(nil), r.Materials...)
	copy.Steps = append([]string(nil), r.Steps...)
	copy.Notes = append([]string(nil), r.Notes...)
	return copy
}

func (r Record) IsVisible() bool { return r.Status == StatusApproved || r.Status == StatusArchived }

func NewRecord(id, title, edition, summary, ageRange string, materials, steps []string, now time.Time) Record {
	return Record{ID: id, Title: title, Edition: edition, Summary: summary, AgeRange: ageRange, Materials: append([]string(nil), materials...), Steps: append([]string(nil), steps...), Status: StatusDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
}
