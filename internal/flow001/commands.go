package flow001

import (
	"example.com/scienceweekly/internal/model"
	"fmt"
	"time"
)

type Command struct {
	Name     string
	RecordID string
	Actor    string
	Content  string
	Approved bool
}
type Result struct {
	Record  model.Record
	Message string
}

func (s *Service) Execute(command Command) (Result, error) {
	switch command.Name {
	case "status":
		record, err := s.AddStatus(command.RecordID, command.Content)
		return Result{Record: record, Message: "状态已记录"}, err
	case "submit":
		record, err := s.SubmitAndDecide(command.RecordID, command.Actor, command.Approved, command.Content)
		return Result{Record: record, Message: "审核流程已更新"}, err
	case "archive":
		record, err := s.ArchiveApproved(command.RecordID, command.Actor)
		return Result{Record: record, Message: "记录已归档"}, err
	default:
		return Result{}, fmt.Errorf("unknown command %q", command.Name)
	}
}

func (s *Service) WithTime(now time.Time) *Service { copy := *s; copy.now = now; return &copy }
