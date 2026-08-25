package registry

import (
	"strings"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Query struct{ store *store.Store }

func NewQuery(s *store.Store) Query { return Query{store: s} }

func (q Query) Search(filter model.SearchFilter) ([]model.Record, error) {
	records, err := q.store.ListRecords(filter)
	if err != nil {
		return nil, err
	}
	return model.SortRecords(records), nil
}

func (q Query) FindVisible(edition string) ([]model.Record, error) {
	records, err := q.Search(model.SearchFilter{Edition: edition, IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	visible := records[:0]
	for _, record := range records {
		if record.IsVisible() {
			visible = append(visible, record)
		}
	}
	return visible, nil
}

func (q Query) Summaries(records []model.Record) []string {
	result := make([]string, 0, len(records))
	for _, record := range records {
		result = append(result, strings.TrimSpace(record.Title)+" - "+strings.TrimSpace(record.Summary))
	}
	return result
}

func (q Query) GroupByEdition(records []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, record := range records {
		groups[record.Edition] = append(groups[record.Edition], record)
	}
	return groups
}
