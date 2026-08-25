package flow001

import (
	"example.com/scienceweekly/internal/model"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	flow, db, _ := newFlowFixture(t)
	defer db.Close()
	record, err := flow.CreateReviewArchive("纸桥", "W01", "承重实验", "9-12", []string{"纸", "胶带"}, []string{"折叠", "加重", "记录"}, "editor", "teacher")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != model.StatusArchived || record.ArchivedAt.IsZero() {
		t.Fatalf("record=%#v", record)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	flow, db, query := newFlowFixture(t)
	defer db.Close()
	if _, err := flow.CreateReviewArchive("磁铁", "W02", "磁场", "6-8", []string{"磁铁"}, []string{"靠近", "观察", "记录"}, "editor", "teacher"); err != nil {
		t.Fatal(err)
	}
	items, err := query.FindVisible("W02")
	if err != nil || len(items) != 1 || items[0].Title != "磁铁" {
		t.Fatalf("items=%v err=%v", items, err)
	}
}
