package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"example.com/scienceweekly/internal/model"
)

type Manifest struct {
	Edition     string
	GeneratedAt time.Time
	RecordIDs   []string
	Digest      string
}

func (r Reporter) Manifest(edition string, now time.Time) (Manifest, error) {
	records, err := r.store.ListRecords(model.SearchFilter{Edition: edition, IncludeArchived: true})
	if err != nil {
		return Manifest{}, err
	}
	ids := make([]string, 0, len(records))
	for _, record := range records {
		if record.Status == model.StatusArchived {
			ids = append(ids, record.ID)
		}
	}
	sort.Strings(ids)
	hash := sha256.Sum256([]byte(strings.Join(ids, "\n")))
	return Manifest{Edition: edition, GeneratedAt: now, RecordIDs: ids, Digest: hex.EncodeToString(hash[:])}, nil
}

func (m Manifest) Contains(recordID string) bool {
	for _, id := range m.RecordIDs {
		if id == recordID {
			return true
		}
	}
	return false
}

func (m Manifest) IsEmpty() bool { return len(m.RecordIDs) == 0 }

func (m Manifest) Summary() string {
	return m.Edition + ":" + string(rune(len(m.RecordIDs))) + ":" + m.Digest
}
