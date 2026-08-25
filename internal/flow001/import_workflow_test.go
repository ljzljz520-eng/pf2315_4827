package flow001

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/archive"
	"example.com/scienceweekly/internal/importer"
	"example.com/scienceweekly/internal/model"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/review"
	"example.com/scienceweekly/internal/store"
)

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/import-flow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	reg := registry.New(db, "batch")
	result := importer.New(reg).Import([]model.ImportRow{{Title: "水滴", Edition: "W05", Summary: "表面张力", AgeRange: "9-12", Materials: "水|硬币", Steps: "滴水|观察|记录"}}, "teacher", time.Unix(0, 0))
	if len(result.Created) != 1 {
		t.Fatalf("result=%#v", result)
	}
	reviewService := review.New(db, "batch")
	record := result.Created[0]
	if _, err := reviewService.Submit(record.ID, "teacher", time.Unix(1, 0)); err != nil {
		t.Fatal(err)
	}
	if _, err := reviewService.Decide(model.ReviewDecision{RecordID: record.ID, Approved: true, Reviewer: "reviewer"}, time.Unix(2, 0)); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.New(db, "batch").Archive(record.ID, "reviewer", time.Unix(3, 0)); err != nil {
		t.Fatal(err)
	}
	report, err := archive.NewReporter(db).Edition("W05", time.Unix(4, 0))
	if err != nil || report.Archived != 1 {
		t.Fatalf("report=%#v err=%v", report, err)
	}
}
