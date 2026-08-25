package store

import (
	"example.com/scienceweekly/internal/model"
	"fmt"
	"go.etcd.io/bbolt"
)

func (s *Store) PutWorkflow(workflow model.Workflow) error {
	data, err := marshal(workflow)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(workflowBucket).Put(key(workflow.ID), data) })
}

func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return model.Workflow{}, err
	}
	var workflow model.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(workflowBucket).Get(key(id))
		if value == nil {
			return fmt.Errorf("workflow %s not found", id)
		}
		return unmarshal(value, &workflow)
	})
	return workflow, err
}

func (s *Store) ListWorkflows(recordID string) ([]model.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	result := make([]model.Workflow, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(workflowBucket)
		for _, k := range sortedKeys(bucket) {
			var workflow model.Workflow
			if err := unmarshal(bucket.Get(k), &workflow); err != nil {
				return err
			}
			if recordID == "" || workflow.RecordID == recordID {
				result = append(result, workflow)
			}
		}
		return nil
	})
	return result, err
}

func (s *Store) PutAttachment(attachment model.Attachment) error {
	data, err := marshal(attachment)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(attachmentBucket).Put(key(attachment.ID), data) })
}

func (s *Store) ListAttachments(recordID string) ([]model.Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	result := make([]model.Attachment, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(attachmentBucket)
		for _, k := range sortedKeys(bucket) {
			var attachment model.Attachment
			if err := unmarshal(bucket.Get(k), &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				result = append(result, attachment)
			}
		}
		return nil
	})
	return result, err
}
