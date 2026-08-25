package flow001

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/archive"
	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/review"
	"example.com/scienceweekly/internal/store"
)

func newFlowFixture(t *testing.T) (*Service, *store.Store, registry.Query) {
	t.Helper()
	db, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	reg := registry.New(db, "flow")
	rev := review.New(db, "flow")
	arch := archive.New(db, "flow")
	return New(reg, rev, arch, now), db, registry.NewQuery(db)
}
