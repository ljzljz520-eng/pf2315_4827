package importer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"example.com/scienceweekly/internal/model"
)

func ReadRows(reader io.Reader) ([]model.ImportRow, error) {
	scanner := bufio.NewScanner(reader)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read rows: %w", err)
	}
	return ParseLines(lines)
}

func WriteRows(rows []model.ImportRow) string {
	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, strings.Join([]string{row.Title, row.Edition, row.Summary, row.AgeRange, row.Materials, row.Steps}, "\t"))
	}
	return strings.Join(lines, "\n")
}

func RowCount(rows []model.ImportRow) int { return len(rows) }

func RejectedByRow(errors []RowError) map[int]string {
	result := make(map[int]string)
	for _, rowError := range errors {
		result[rowError.Row] = rowError.Message
	}
	return result
}
