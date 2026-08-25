package store

import (
	"fmt"

	"example.com/scienceweekly/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) PutAudit(event model.AuditEvent) error {
	data, err := marshal(event)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put(key(event.ID), data) })
}

func (s *Store) ListAudits(recordID string) ([]model.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	result := make([]model.AuditEvent, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(auditBucket)
		for _, k := range sortedKeys(bucket) {
			var event model.AuditEvent
			if err := unmarshal(bucket.Get(k), &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				result = append(result, event)
			}
		}
		return nil
	})
	return result, err
}

func (s *Store) RequireAudit(id string) (model.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return model.AuditEvent{}, err
	}
	var event model.AuditEvent
	err := s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(auditBucket).Get(key(id))
		if value == nil {
			return fmt.Errorf("audit %s not found", id)
		}
		return unmarshal(value, &event)
	})
	return event, err
}
