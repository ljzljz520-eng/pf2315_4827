package model

import (
	"fmt"
	"strings"
	"time"
)

func NewWorkflow(id, recordID, name, owner string, steps []string, started time.Time) Workflow {
	return Workflow{ID: id, RecordID: recordID, Name: strings.TrimSpace(name), Owner: strings.TrimSpace(owner), State: "active", CurrentStep: 0, Steps: normalizeList(steps), StartedAt: started}
}

func (w Workflow) Validate() error {
	if strings.TrimSpace(w.ID) == "" {
		return fmt.Errorf("workflow id is required")
	}
	if strings.TrimSpace(w.RecordID) == "" {
		return fmt.Errorf("workflow record id is required")
	}
	if strings.TrimSpace(w.Name) == "" {
		return fmt.Errorf("workflow name is required")
	}
	if len(w.Steps) == 0 {
		return fmt.Errorf("workflow must contain steps")
	}
	if w.CurrentStep < 0 || w.CurrentStep > len(w.Steps) {
		return fmt.Errorf("workflow step %d is outside range", w.CurrentStep)
	}
	if w.State != "active" && w.State != "completed" && w.State != "cancelled" {
		return fmt.Errorf("unknown workflow state %q", w.State)
	}
	return nil
}

func (w *Workflow) Advance(now time.Time) error {
	if err := w.Validate(); err != nil {
		return err
	}
	if w.State != "active" {
		return fmt.Errorf("workflow is %s", w.State)
	}
	if w.CurrentStep >= len(w.Steps) {
		return fmt.Errorf("workflow is already complete")
	}
	w.CurrentStep++
	if w.CurrentStep == len(w.Steps) {
		w.State = "completed"
		w.FinishedAt = now
	}
	return nil
}

func (w *Workflow) Cancel(now time.Time) error {
	if err := w.Validate(); err != nil {
		return err
	}
	if w.State != "active" {
		return fmt.Errorf("workflow is %s", w.State)
	}
	w.State = "cancelled"
	w.FinishedAt = now
	return nil
}

func (w Workflow) CurrentAction() string {
	if w.CurrentStep >= len(w.Steps) {
		return "完成"
	}
	return w.Steps[w.CurrentStep]
}

func (a Attachment) Validate() error {
	if strings.TrimSpace(a.ID) == "" {
		return fmt.Errorf("attachment id is required")
	}
	if strings.TrimSpace(a.RecordID) == "" {
		return fmt.Errorf("attachment record id is required")
	}
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("attachment name is required")
	}
	if len(a.Content) == 0 {
		return fmt.Errorf("attachment content is required")
	}
	return nil
}

func (a Attachment) IsImage() bool { return strings.HasPrefix(strings.ToLower(a.MediaType), "image/") }
