package store

import (
	"livestream/internal/model"
)

func (s *MemoryStore) CreateDanmaku(d *model.Danmaku) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.danmakus[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDanmaku(id string) (*model.Danmaku, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.danmakus[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDanmakus() []*model.Danmaku {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Danmaku, 0, len(s.danmakus))
	for _, d := range s.danmakus {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDanmaku(d *model.Danmaku) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.danmakus[d.ID]; !ok {
		return ErrNotFound
	}
	s.danmakus[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDanmaku(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.danmakus[id]; !ok {
		return ErrNotFound
	}
	delete(s.danmakus, id)
	return nil
}
