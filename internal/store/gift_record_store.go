package store

import (
	"livestream/internal/model"
)

func (s *MemoryStore) CreateGiftRecord(r *model.GiftRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.giftRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) GetGiftRecord(id string) (*model.GiftRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.giftRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListGiftRecords() []*model.GiftRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.GiftRecord, 0, len(s.giftRecords))
	for _, r := range s.giftRecords {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) DeleteGiftRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.giftRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.giftRecords, id)
	return nil
}
