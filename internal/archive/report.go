package archive

import (
	"sort"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Reporter struct{ store *store.Store }

func NewReporter(s *store.Store) Reporter { return Reporter{store: s} }

func (r Reporter) Edition(edition string, now time.Time) (model.EditionReport, error) {
	records, err := r.store.ListRecords(model.SearchFilter{Edition: edition, IncludeArchived: true})
	if err != nil {
		return model.EditionReport{}, err
	}
	report := model.EditionReport{Edition: edition, Total: len(records), GeneratedAt: now, VisibleTitles: make([]string, 0)}
	for _, record := range records {
		switch record.Status {
		case model.StatusDraft:
			report.Draft++
		case model.StatusSubmitted:
			report.Submitted++
		case model.StatusApproved:
			report.Approved++
		case model.StatusArchived:
			report.Archived++
		case model.StatusRejected:
			report.Rejected++
		}
		if record.IsVisible() {
			report.VisibleTitles = append(report.VisibleTitles, record.Title)
		}
	}
	sort.Strings(report.VisibleTitles)
	return report, nil
}

func (r Reporter) AuditSummary(recordID string) (map[string]int, error) {
	events, err := r.store.ListAudits(recordID)
	if err != nil {
		return nil, err
	}
	result := map[string]int{}
	for _, event := range events {
		result[event.Action]++
	}
	return result, nil
}
