package flow001

import "testing"

func Test2315BusinessRegression(t *testing.T) {
	flow, db, _ := newFlowFixture(t)
	defer db.Close()
	record, err := flow.CreateReviewArchive("状态实验", "W23", "观察状态变化", "9-12", []string{"水"}, []string{"准备", "观察", "记录"}, "editor", "teacher")
	if err != nil {
		t.Fatal(err)
	}
	first, err := flow.AddStatus(record.ID, "第一条状态")
	if err != nil {
		t.Fatal(err)
	}
	second, err := flow.AddStatus(record.ID, "第二条状态")
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Notes) != 1 || first.Notes[0] != "第一条状态" {
		t.Fatalf("first=%v", first.Notes)
	}
	if len(second.Notes) != 2 || second.Notes[1] != "第二条状态" {
		t.Fatalf("second=%v", second.Notes)
	}
}
