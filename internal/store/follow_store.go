package store

import (
	"livestream/internal/model"
)

func (s *MemoryStore) CreateFollow(f *model.Follow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.follows {
		if exist.FollowerID == f.FollowerID && exist.FolloweeID == f.FolloweeID {
			return ErrConflict
		}
	}
	s.follows[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFollow(id string) (*model.Follow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.follows[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) ListFollows() []*model.Follow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Follow, 0, len(s.follows))
	for _, f := range s.follows {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) DeleteFollow(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.follows[id]; !ok {
		return ErrNotFound
	}
	delete(s.follows, id)
	return nil
}

func (s *MemoryStore) FollowExists(followerID, followeeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.follows {
		if f.FollowerID == followerID && f.FolloweeID == followeeID {
			return true
		}
	}
	return false
}
