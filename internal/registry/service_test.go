package registry

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/store"
)

func TestRegisterAndAmend(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/registry.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db, "science")
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	record, err := svc.Register("彩虹密度", "W03", "比较液体", "9-12", []string{"水", "盐"}, []string{"装杯", "观察"}, "小明", now)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := svc.Amend(record.ID, "小红", "彩虹密度实验", "比较不同液体", nil, []string{"装杯", "加盐", "观察"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 2 || updated.Editor != "小红" || updated.Title != "彩虹密度实验" {
		t.Fatalf("updated=%#v", updated)
	}
}
