package importer

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
)

type Importer struct{ registry *registry.Service }

func New(r *registry.Service) *Importer { return &Importer{registry: r} }

type Result struct {
	Created  []model.Record
	Rejected []RowError
}
type RowError struct {
	Row     int
	Message string
}

func (i *Importer) Import(rows []model.ImportRow, editor string, now time.Time) Result {
	result := Result{Created: make([]model.Record, 0, len(rows)), Rejected: make([]RowError, 0)}
	for index, row := range rows {
		materials := split(row.Materials)
		steps := split(row.Steps)
		record, err := i.registry.Register(row.Title, row.Edition, row.Summary, row.AgeRange, materials, steps, editor, now)
		if err != nil {
			result.Rejected = append(result.Rejected, RowError{Row: index + 1, Message: err.Error()})
			continue
		}
		result.Created = append(result.Created, record)
	}
	return result
}

func split(value string) []string {
	parts := strings.Split(value, "|")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ParseLines(lines []string) ([]model.ImportRow, error) {
	rows := make([]model.ImportRow, 0, len(lines))
	for index, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 6 {
			return nil, fmt.Errorf("line %d has %d columns", index+1, len(fields))
		}
		rows = append(rows, model.ImportRow{Title: fields[0], Edition: fields[1], Summary: fields[2], AgeRange: fields[3], Materials: fields[4], Steps: fields[5]})
	}
	return rows, nil
}

func FormatResult(result Result) string {
	return fmt.Sprintf("created=%d rejected=%d", len(result.Created), len(result.Rejected))
}
