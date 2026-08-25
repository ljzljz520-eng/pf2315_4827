package model

import (
	"fmt"
	"strings"
)

func ParseStatus(raw string) (Status, error) {
	switch Status(strings.ToLower(strings.TrimSpace(raw))) {
	case StatusDraft:
		return StatusDraft, nil
	case StatusSubmitted:
		return StatusSubmitted, nil
	case StatusApproved:
		return StatusApproved, nil
	case StatusArchived:
		return StatusArchived, nil
	case StatusRejected:
		return StatusRejected, nil
	default:
		return "", fmt.Errorf("unknown status %q", raw)
	}
}

func StatusLabel(status Status) string {
	switch status {
	case StatusDraft:
		return "草稿"
	case StatusSubmitted:
		return "待审核"
	case StatusApproved:
		return "已批准"
	case StatusArchived:
		return "已归档"
	case StatusRejected:
		return "已退回"
	default:
		return "未知"
	}
}

func IsTerminal(status Status) bool { return status == StatusArchived }

func StatusOrder(status Status) int {
	order := map[Status]int{StatusDraft: 1, StatusSubmitted: 2, StatusRejected: 2, StatusApproved: 3, StatusArchived: 4}
	return order[status]
}
