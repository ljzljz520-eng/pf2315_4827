package flow001

import (
	"fmt"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
)

func (s *Service) StartEditorialWorkflow(recordID, owner string, steps []string) (model.Workflow, error) {
	record, err := s.registry.Load(recordID)
	if err != nil {
		return model.Workflow{}, err
	}
	if len(steps) < 4 {
		return model.Workflow{}, fmt.Errorf("editorial workflow requires four steps")
	}
	workflowID := recordID + "-editorial"
	workflow := model.NewWorkflow(workflowID, record.ID, "editorial", owner, steps, s.now)
	if err := s.store.SaveWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (s *Service) AdvanceEditorialWorkflow(workflowID string) (model.Workflow, error) {
	return s.store.AdvanceWorkflow(workflowID, s.now)
}

func (s *Service) FinishEditorialWorkflow(workflowID string) (model.Workflow, error) {
	workflow, err := s.store.GetWorkflow(workflowID)
	if err != nil {
		return model.Workflow{}, err
	}
	for workflow.State == "active" {
		workflow, err = s.store.AdvanceWorkflow(workflow.ID, s.now)
		if err != nil {
			return model.Workflow{}, err
		}
	}
	return workflow, nil
}

func (s *Service) CancelEditorialWorkflow(workflowID string) (model.Workflow, error) {
	workflow, err := s.store.GetWorkflow(workflowID)
	if err != nil {
		return model.Workflow{}, err
	}
	if err := workflow.Cancel(s.now); err != nil {
		return model.Workflow{}, err
	}
	if err := s.store.PutWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (s *Service) WorkflowProgress(workflowID string) (string, error) {
	workflow, err := s.store.GetWorkflow(workflowID)
	if err != nil {
		return "", err
	}
	if workflow.State == "completed" {
		return "completed", nil
	}
	if workflow.State == "cancelled" {
		return "cancelled", nil
	}
	return fmt.Sprintf("%d/%d:%s", workflow.CurrentStep, len(workflow.Steps), strings.TrimSpace(workflow.CurrentAction())), nil
}

func (s *Service) WorkflowForRecord(recordID string) ([]model.Workflow, error) {
	return s.store.ListWorkflows(recordID)
}

func (s *Service) AttachEvidence(recordID, name, mediaType string, content []byte, checksum string) (model.Attachment, error) {
	return s.store.Attach(recordID, name, mediaType, content, checksum, s.now)
}

func (s *Service) Integrity() (string, error) {
	report, err := s.store.CheckIntegrity()
	if err != nil {
		return "", err
	}
	if len(report.Errors) > 0 {
		return strings.Join(report.Errors, ";"), nil
	}
	return fmt.Sprintf("records=%d audits=%d workflows=%d attachments=%d", report.Records, report.Audits, report.Workflows, report.Attachments), nil
}

func (s *Service) ArchiveBefore(cutoff time.Time) (int, error) {
	return s.store.DeleteArchivedBefore(cutoff)
}
