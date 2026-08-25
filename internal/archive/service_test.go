package archive

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/review"
	"example.com/scienceweekly/internal/store"
)

func TestArchiveEditionAndReport(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/archive.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Unix(10, 0)
	reg := registry.New(db, "a")
	rev := review.New(db, "v")
	svc := New(db, "x")
	record, err := reg.Register("影子", "W04", "光线", "6-8", []string{"灯"}, []string{"照射", "观察"}, "e", now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := rev.Submit(record.ID, "e", now); err != nil {
		t.Fatal(err)
	}
	if _, err := rev.Decide(model.ReviewDecision{RecordID: record.ID, Approved: true, Reviewer: "r"}, now); err != nil {
		t.Fatal(err)
	}
	batch, err := svc.ArchiveEdition("W04", "r", now)
	if err != nil || len(batch.RecordIDs) != 1 {
		t.Fatalf("batch=%v err=%v", batch, err)
	}
	report, err := NewReporter(db).Edition("W04", now)
	if err != nil || report.Archived != 1 || !report.HasContent() {
		t.Fatalf("report=%#v err=%v", report, err)
	}
}
