package importer

import (
	"testing"
	"time"

	"example.com/scienceweekly/internal/registry"
	"example.com/scienceweekly/internal/store"
)

func TestImportRowsAndParseLines(t *testing.T) {
	rows, err := ParseLines([]string{"纸桥\tW01\t承重\t9-12\t纸|胶带\t折叠|测试|记录"})
	if err != nil || len(rows) != 1 {
		t.Fatalf("rows=%v err=%v", rows, err)
	}
	db, err := store.Open(t.TempDir() + "/import.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	result := New(registry.New(db, "i")).Import(rows, "importer", time.Unix(0, 0))
	if len(result.Created) != 1 || len(result.Rejected) != 0 {
		t.Fatalf("result=%#v", result)
	}
	if FormatResult(result) != "created=1 rejected=0" {
		t.Fatalf("format=%s", FormatResult(result))
	}
}
