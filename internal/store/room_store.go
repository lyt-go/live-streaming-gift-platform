package store

import (
	"livestream/internal/model"
)

func (s *MemoryStore) CreateRoom(r *model.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRoom(id string) (*model.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rooms[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRooms() []*model.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Room, 0, len(s.rooms))
	for _, r := range s.rooms {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRoom(r *model.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[r.ID]; !ok {
		return ErrNotFound
	}
	s.rooms[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRoom(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[id]; !ok {
		return ErrNotFound
	}
	delete(s.rooms, id)
	return nil
}
