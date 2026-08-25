package store

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/model"
)

func TestStoreSearchFilters(t *testing.T) {
	db, err := Open(t.TempDir() + "/filter.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, record := range []model.Record{
		model.NewRecord("a", "磁铁", "W01", "磁场", "6-8", []string{"磁铁"}, []string{"靠近", "观察"}, now),
		model.NewRecord("b", "植物", "W02", "发芽", "9-12", []string{"种子"}, []string{"浇水", "观察"}, now),
	} {
		if err := db.PutRecord(record); err != nil {
			t.Fatal(err)
		}
	}
	items, err := db.ListRecords(model.SearchFilter{Text: "磁"})
	if err != nil || len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("items=%v err=%v", items, err)
	}
	counts, err := db.SnapshotCounts()
	if err != nil || counts["records"] != 2 {
		t.Fatalf("counts=%v err=%v", counts, err)
	}
}
