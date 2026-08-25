package store

import (
	"fmt"
	"strings"

	"example.com/scienceweekly/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) PutRecord(record model.Record) error {
	data, err := marshal(record)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(recordBucket).Put(key(record.ID), data) })
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return model.Record{}, err
	}
	var record model.Record
	err := s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(recordBucket).Get(key(id))
		if value == nil {
			return fmt.Errorf("record %s not found", id)
		}
		return unmarshal(value, &record)
	})
	return record, err
}

func (s *Store) DeleteRecord(id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(recordBucket).Delete(key(id)) })
}

func (s *Store) ListRecords(filter model.SearchFilter) ([]model.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	result := make([]model.Record, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(recordBucket)
		for _, k := range sortedKeys(bucket) {
			var record model.Record
			if err := unmarshal(bucket.Get(k), &record); err != nil {
				return err
			}
			if matches(record, filter) {
				result = append(result, record)
			}
		}
		return nil
	})
	return result, err
}

func matches(record model.Record, filter model.SearchFilter) bool {
	if !filter.IncludeArchived && record.Status == model.StatusArchived {
		return false
	}
	if filter.Status != "" && record.Status != filter.Status {
		return false
	}
	if filter.Edition != "" && !strings.EqualFold(record.Edition, filter.Edition) {
		return false
	}
	if filter.AgeRange != "" && !strings.EqualFold(record.AgeRange, filter.AgeRange) {
		return false
	}
	if filter.Text != "" {
		text := strings.ToLower(filter.Text)
		if !strings.Contains(strings.ToLower(record.Title), text) && !strings.Contains(strings.ToLower(record.Summary), text) {
			return false
		}
	}
	return true
}

func (s *Store) CountRecords(filter model.SearchFilter) (int, error) {
	records, err := s.ListRecords(filter)
	return len(records), err
}
