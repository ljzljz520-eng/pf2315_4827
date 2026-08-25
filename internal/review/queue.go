package review

import (
	"sort"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/store"
)

type Queue struct{ store *store.Store }

func NewQueue(s *store.Store) Queue { return Queue{store: s} }

func (q Queue) Pending(edition string) ([]model.Record, error) {
	filter := model.SearchFilter{Status: model.StatusSubmitted, Edition: edition}
	items, err := q.store.ListRecords(filter)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].UpdatedAt.Before(items[j].UpdatedAt)
	})
	return items, nil
}

func (q Queue) Assign(record model.Record, reviewer string, now time.Time) (model.Record, error) {
	reviewer = strings.TrimSpace(reviewer)
	if reviewer == "" {
		return model.Record{}, &QueueError{"reviewer is required"}
	}
	if record.Status != model.StatusSubmitted {
		return model.Record{}, &QueueError{"record is not waiting for review"}
	}
	record.Reviewer = reviewer
	record.UpdatedAt = now
	if err := q.store.PutRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

type QueueError struct{ Message string }

func (e *QueueError) Error() string { return e.Message }

func (q Queue) IsOverdue(record model.Record, cutoff time.Time) bool {
	return record.Status == model.StatusSubmitted && record.UpdatedAt.Before(cutoff)
}

func (q Queue) GroupByReviewer(records []model.Record) map[string][]model.Record {
	groups := make(map[string][]model.Record)
	for _, record := range records {
		key := record.Reviewer
		if key == "" {
			key = "unassigned"
		}
		groups[key] = append(groups[key], record)
	}
	return groups
}

func (q Queue) Workload(records []model.Record, reviewer string) int {
	count := 0
	for _, record := range records {
		if record.Reviewer == reviewer && record.Status == model.StatusSubmitted {
			count++
		}
	}
	return count
}
