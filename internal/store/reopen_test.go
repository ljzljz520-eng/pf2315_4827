package store

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/weekly.db"
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	record := model.NewRecord("persist-1", "纸桥承重", "2026-W01", "测试纸张结构", "9-12", []string{"纸"}, []string{"折叠", "测量", "记录"}, now)
	if err := first.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.PutAudit(model.AuditEvent{ID: "audit-1", RecordID: record.ID, Action: "create", CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Title != record.Title || loaded.Summary != record.Summary {
		t.Fatalf("loaded record = %#v", loaded)
	}
	events, err := second.ListAudits(record.ID)
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%v err=%v", events, err)
	}
}
