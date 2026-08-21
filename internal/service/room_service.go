package service

import (
	"sort"
	"time"

	"livestream/internal/model"
	"livestream/internal/store"
	"livestream/pkg/idgen"
)

// CreateRoom 创建直播间，主播必须存在。
func (s *Service) CreateRoom(input model.Room) (*model.Room, error) {
	if _, err := s.store.GetUser(input.StreamerID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("streamer_id", "主播不存在")
		}
		return nil, err
	}
	r := &model.Room{
		ID:          idgen.Hex(),
		Title:       input.Title,
		Category:    input.Category,
		StreamerID:  input.StreamerID,
		CoverURL:    input.CoverURL,
		Status:      input.Status,
		ViewerCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateRoom(r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetRoom 按 ID 查询直播间。
func (s *Service) GetRoom(id string) (*model.Room, error) {
	return s.store.GetRoom(id)
}

// ListRooms 分页查询直播间列表。
func (s *Service) ListRooms(filter model.RoomFilter, page, size int) ([]*model.Room, int, error) {
	all := s.store.ListRooms()
	matched := make([]*model.Room, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].ViewerCount != matched[j].ViewerCount {
			return matched[i].ViewerCount > matched[j].ViewerCount
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Room{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateRoom 更新直播间基础信息（标题、分类、封面）。
func (s *Service) UpdateRoom(id string, input model.Room) (*model.Room, error) {
	existing, err := s.store.GetRoom(id)
	if err != nil {
		return nil, err
	}
	existing.Title = input.Title
	existing.Category = input.Category
	existing.CoverURL = input.CoverURL
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRoom(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteRoom 删除直播间。
func (s *Service) DeleteRoom(id string) error {
	return s.store.DeleteRoom(id)
}

// StartRoom 开播：pending -> live。
func (s *Service) StartRoom(id string) (*model.Room, error) {
	r, err := s.store.GetRoom(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(r.Status, model.RoomStatusLive) {
		return nil, model.NewValidationError("status", "当前状态不允许开播")
	}
	r.Status = model.RoomStatusLive
	r.StartedAt = time.Now()
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRoom(r); err != nil {
		return nil, err
	}
	return r, nil
}

// EndRoom 下播：live -> ended。
func (s *Service) EndRoom(id string) (*model.Room, error) {
	r, err := s.store.GetRoom(id)
	if err != nil {
		return nil, err
	}
	r.Status = model.RoomStatusEnded
	if !model.CanTransition(r.Status, model.RoomStatusEnded) {
		return nil, model.NewValidationError("status", "当前状态不允许下播")
	}
	r.Status = model.RoomStatusEnded
	r.EndedAt = time.Now()
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRoom(r); err != nil {
		return nil, err
	}
	return r, nil
}
