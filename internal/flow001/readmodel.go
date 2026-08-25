package flow001

import (
	"sort"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
)

type Dashboard struct {
	Edition     string
	Total       int
	Visible     int
	Draft       int
	Review      int
	Archived    int
	Titles      []string
	GeneratedAt time.Time
}

func (s *Service) Dashboard(edition string, query registry.Query) (Dashboard, error) {
	items, err := query.Search(model.SearchFilter{Edition: edition, IncludeArchived: true})
	if err != nil {
		return Dashboard{}, err
	}
	dashboard := Dashboard{Edition: edition, Total: len(items), Titles: make([]string, 0), GeneratedAt: s.now}
	for _, record := range items {
		switch record.Status {
		case model.StatusDraft, model.StatusRejected:
			dashboard.Draft++
		case model.StatusSubmitted:
			dashboard.Review++
		case model.StatusApproved, model.StatusArchived:
			dashboard.Visible++
			if record.Status == model.StatusArchived {
				dashboard.Archived++
			}
		}
		if record.IsVisible() {
			dashboard.Titles = append(dashboard.Titles, record.Title)
		}
	}
	sort.Strings(dashboard.Titles)
	return dashboard, nil
}

func (d Dashboard) CompletionPercent() int {
	if d.Total == 0 {
		return 0
	}
	return d.Visible * 100 / d.Total
}

func (d Dashboard) ContainsTitle(title string) bool {
	for _, candidate := range d.Titles {
		if strings.EqualFold(candidate, title) {
			return true
		}
	}
	return false
}

func (s *Service) Timeline(record model.Record) []string {
	result := make([]string, 0, len(record.Notes))
	for index, note := range record.Notes {
		result = append(result, strings.Join([]string{record.ID, string(record.Status), string(rune(index + 1)), note}, ":"))
	}
	return result
}

func (s *Service) StatusCount(record model.Record, content string) int {
	count := 0
	for _, note := range record.Notes {
		if note == content {
			count++
		}
	}
	return count
}
