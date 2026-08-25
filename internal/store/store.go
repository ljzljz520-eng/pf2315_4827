package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var (
	recordBucket     = []byte("records")
	auditBucket      = []byte("audits")
	workflowBucket   = []byte("workflows")
	attachmentBucket = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db}
	if err := db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{recordBucket, auditBucket, workflowBucket, attachmentBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize buckets: %w", err)
	}
	return s, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) ensureOpen() error {
	if s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}

func marshal(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode value: %w", err)
	}
	return data, nil
}

func unmarshal(data []byte, value any) error {
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}
	return nil
}

func key(value string) []byte { return []byte(value) }

func sortedKeys(bucket *bbolt.Bucket) [][]byte {
	keys := make([][]byte, 0)
	_ = bucket.ForEach(func(k, v []byte) error {
		if v != nil {
			keys = append(keys, append([]byte(nil), k...))
		}
		return nil
	})
	sort.Slice(keys, func(i, j int) bool { return string(keys[i]) < string(keys[j]) })
	return keys
}

func (s *Store) SnapshotCounts() (map[string]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	counts := map[string]int{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for name, bucket := range map[string][]byte{"records": recordBucket, "audits": auditBucket, "workflows": workflowBucket, "attachments": attachmentBucket} {
			counts[name] = tx.Bucket(bucket).Stats().KeyN
		}
		return nil
	})
	return counts, err
}
