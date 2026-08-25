package model

import "time"

type EditionReport struct {
	Edition       string
	Total         int
	Draft         int
	Submitted     int
	Approved      int
	Archived      int
	Rejected      int
	VisibleTitles []string
	GeneratedAt   time.Time
}

type ReviewDecision struct {
	RecordID string
	Approved bool
	Reviewer string
	Reason   string
}

type ArchiveBatch struct {
	Edition    string
	RecordIDs  []string
	Actor      string
	ArchivedAt time.Time
}

func (r EditionReport) CompletionRate() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Approved+r.Archived) / float64(r.Total)
}

func (r EditionReport) HasContent() bool { return len(r.VisibleTitles) > 0 }
