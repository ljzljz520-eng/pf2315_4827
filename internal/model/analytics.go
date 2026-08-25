package model

import "strings"

type StatusCounts struct {
	Draft     int
	Submitted int
	Approved  int
	Rejected  int
	Archived  int
}

func CountStatuses(records []Record) StatusCounts {
	counts := StatusCounts{}
	for _, record := range records {
		switch record.Status {
		case StatusDraft:
			counts.Draft++
		case StatusSubmitted:
			counts.Submitted++
		case StatusApproved:
			counts.Approved++
		case StatusRejected:
			counts.Rejected++
		case StatusArchived:
			counts.Archived++
		}
	}
	return counts
}

func (c StatusCounts) Total() int {
	return c.Draft + c.Submitted + c.Approved + c.Rejected + c.Archived
}

func (c StatusCounts) Published() int { return c.Approved + c.Archived }

func (c StatusCounts) IsBalanced() bool {
	return c.Total() == c.Draft+c.Submitted+c.Approved+c.Rejected+c.Archived
}

func JoinMaterials(records []Record) string {
	values := make([]string, 0)
	seen := make(map[string]bool)
	for _, record := range records {
		for _, material := range record.Materials {
			material = strings.TrimSpace(material)
			if material != "" && !seen[material] {
				seen[material] = true
				values = append(values, material)
			}
		}
	}
	return strings.Join(values, ", ")
}

func LongestStep(records []Record) string {
	longest := ""
	for _, record := range records {
		for _, step := range record.Steps {
			if len(step) > len(longest) {
				longest = step
			}
		}
	}
	return longest
}
