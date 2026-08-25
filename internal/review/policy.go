package review

import (
	"example.com/scienceweekly/internal/model"
	"strings"
)

type Policy struct {
	RequiredAgeRanges []string
	MinimumSteps      int
	RestrictedWords   []string
}

func DefaultPolicy() Policy {
	return Policy{RequiredAgeRanges: []string{"6-8", "9-12"}, MinimumSteps: 3, RestrictedWords: []string{"火焰", "强酸", "爆炸"}}
}

func (p Policy) Check(record model.Record) []model.ValidationIssue {
	issues := make([]model.ValidationIssue, 0)
	if p.MinimumSteps > 0 && len(record.Steps) < p.MinimumSteps {
		issues = append(issues, model.ValidationIssue{"steps", "实验步骤不足"})
	}
	if len(p.RequiredAgeRanges) > 0 {
		matched := false
		for _, age := range p.RequiredAgeRanges {
			if strings.EqualFold(age, record.AgeRange) {
				matched = true
			}
		}
		if !matched {
			issues = append(issues, model.ValidationIssue{"age_range", "年龄范围不在周刊范围"})
		}
	}
	text := strings.ToLower(record.Title + " " + record.Summary + " " + strings.Join(record.Steps, " "))
	for _, word := range p.RestrictedWords {
		if strings.Contains(text, strings.ToLower(word)) {
			issues = append(issues, model.ValidationIssue{"content", "包含限制内容"})
		}
	}
	return issues
}

func (p Policy) Allows(record model.Record) bool { return len(p.Check(record)) == 0 }

func (p Policy) Explain(record model.Record) string {
	issues := p.Check(record)
	if len(issues) == 0 {
		return "通过"
	}
	parts := make([]string, len(issues))
	for i, issue := range issues {
		parts[i] = issue.Field + ":" + issue.Message
	}
	return strings.Join(parts, "; ")
}
