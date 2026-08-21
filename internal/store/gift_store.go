package store

import (
	"livestream/internal/model"
)

func (s *MemoryStore) CreateGift(g *model.Gift) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.gifts {
		if exist.Name == g.Name {
			return ErrConflict
		}
	}
	s.gifts[g.ID] = g
	return nil
}

func (s *MemoryStore) GetGift(id string) (*model.Gift, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.gifts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return g, nil
}

func (s *MemoryStore) ListGifts() []*model.Gift {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Gift, 0, len(s.gifts))
	for _, g := range s.gifts {
		list = append(list, g)
	}
	return list
}

func (s *MemoryStore) UpdateGift(g *model.Gift) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.gifts[g.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.gifts {
		if exist.ID != g.ID && exist.Name == g.Name {
			return ErrConflict
		}
	}
	s.gifts[g.ID] = g
	return nil
}

func (s *MemoryStore) DeleteGift(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.gifts[id]; !ok {
		return ErrNotFound
	}
	delete(s.gifts, id)
	return nil
}
