package review

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/store"
)

func TestReviewApprovalCreatesAudit(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/review.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	registryService := registry.New(db, "r")
	record, err := registryService.Register("种子", "W01", "发芽观察", "6-8", []string{"种子"}, []string{"浸泡", "观察"}, "editor", time.Unix(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	svc := New(db, "review")
	if _, err := svc.Submit(record.ID, "editor", time.Unix(1, 0)); err != nil {
		t.Fatal(err)
	}
	approved, err := svc.Decide(model.ReviewDecision{RecordID: record.ID, Approved: true, Reviewer: "teacher"}, time.Unix(2, 0))
	if err != nil {
		t.Fatal(err)
	}
	if approved.Status != model.StatusApproved {
		t.Fatalf("status=%s", approved.Status)
	}
	history, err := svc.History(record.ID)
	if err != nil || len(history) != 2 {
		t.Fatalf("history=%v err=%v", history, err)
	}
}
