package model

import (
	"fmt"
	"sort"
	"strings"
)

func NormalizeRecord(record Record) Record {
	record.Title = strings.TrimSpace(record.Title)
	record.Edition = strings.TrimSpace(record.Edition)
	record.Summary = strings.TrimSpace(record.Summary)
	record.AgeRange = strings.TrimSpace(record.AgeRange)
	record.Editor = strings.TrimSpace(record.Editor)
	record.Reviewer = strings.TrimSpace(record.Reviewer)
	record.Materials = normalizeList(record.Materials)
	record.Steps = normalizeList(record.Steps)
	record.Notes = normalizeList(record.Notes)
	return record
}

func normalizeList(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]bool)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" && !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}

func SortRecords(records []Record) []Record {
	result := make([]Record, len(records))
	copy(result, records)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Edition != result[j].Edition {
			return result[i].Edition < result[j].Edition
		}
		if StatusOrder(result[i].Status) != StatusOrder(result[j].Status) {
			return StatusOrder(result[i].Status) > StatusOrder(result[j].Status)
		}
		return result[i].ID < result[j].ID
	})
	return result
}

func RecordHeadline(record Record) string {
	return fmt.Sprintf("%s [%s] %s", record.Title, StatusLabel(record.Status), record.Edition)
}

func RecordDigest(record Record) string {
	parts := []string{record.ID, record.Title, record.Edition, record.Summary, string(record.Status)}
	return strings.Join(parts, "|")
}

func HasMaterial(record Record, material string) bool {
	material = strings.TrimSpace(material)
	for _, value := range record.Materials {
		if strings.EqualFold(value, material) {
			return true
		}
	}
	return false
}

func StepNumber(record Record, step string) int {
	for index, value := range record.Steps {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(step)) {
			return index + 1
		}
	}
	return 0
}

func MergeNotes(record Record, additions ...string) Record {
	for _, note := range additions {
		note = strings.TrimSpace(note)
		if note != "" {
			record.Notes = append(record.Notes, note)
		}
	}
	return record
}
