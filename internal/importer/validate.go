package importer

import (
	"example.com/scienceweekly/internal/model"
	"strings"
)

func ValidateRows(rows []model.ImportRow) []RowError {
	errors := make([]RowError, 0)
	for index, row := range rows {
		if strings.TrimSpace(row.Title) == "" {
			errors = append(errors, RowError{Row: index + 1, Message: "title is required"})
			continue
		}
		if strings.TrimSpace(row.Edition) == "" {
			errors = append(errors, RowError{Row: index + 1, Message: "edition is required"})
		}
		if len(split(row.Steps)) < 2 {
			errors = append(errors, RowError{Row: index + 1, Message: "steps are required"})
		}
	}
	return errors
}

func NormalizeRows(rows []model.ImportRow) []model.ImportRow {
	result := make([]model.ImportRow, len(rows))
	for i, row := range rows {
		row.Title = strings.TrimSpace(row.Title)
		row.Edition = strings.TrimSpace(row.Edition)
		row.Summary = strings.TrimSpace(row.Summary)
		row.AgeRange = strings.TrimSpace(row.AgeRange)
		row.Materials = strings.TrimSpace(row.Materials)
		row.Steps = strings.TrimSpace(row.Steps)
		result[i] = row
	}
	return result
}
