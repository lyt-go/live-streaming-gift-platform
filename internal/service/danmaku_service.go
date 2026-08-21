package service

import (
	"sort"
	"time"

	"livestream/internal/model"
	"livestream/internal/store"
	"livestream/pkg/idgen"
)

// CreateDanmaku 发送弹幕，直播间必须存在。
func (s *Service) CreateDanmaku(input model.Danmaku) (*model.Danmaku, error) {
	if _, err := s.store.GetRoom(input.RoomID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("room_id", "直播间不存在")
		}
		return nil, err
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("user_id", "用户不存在")
		}
		return nil, err
	}
	d := &model.Danmaku{
		ID:        idgen.Hex(),
		RoomID:    input.RoomID,
		UserID:    input.UserID,
		Content:   input.Content,
		Status:    model.DanmakuStatusPending,
		CreatedAt: time.Now(),
	}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateDanmaku(d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDanmaku 按 ID 查询弹幕。
func (s *Service) GetDanmaku(id string) (*model.Danmaku, error) {
	return s.store.GetDanmaku(id)
}

// ListDanmakus 分页查询弹幕列表。
func (s *Service) ListDanmakus(filter model.DanmakuFilter, page, size int) ([]*model.Danmaku, int, error) {
	all := s.store.ListDanmakus()
	matched := make([]*model.Danmaku, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Danmaku{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// transitionDanmaku 执行弹幕状态流转。
func (s *Service) transitionDanmaku(id, to string) (*model.Danmaku, error) {
	d, err := s.store.GetDanmaku(id)
	if err != nil {
		return nil, err
	}
	d.Status = to
	if !model.CanTransitionDanmaku(d.Status, to) {
		return nil, model.NewValidationError("status", "当前状态不允许该审核操作")
	}
	d.Status = to
	if err := s.store.UpdateDanmaku(d); err != nil {
		return nil, err
	}
	return d, nil
}

// ApproveDanmaku 审核通过弹幕。
func (s *Service) ApproveDanmaku(id string) (*model.Danmaku, error) {
	return s.transitionDanmaku(id, model.DanmakuStatusApproved)
}

// BlockDanmaku 拦截弹幕。
func (s *Service) BlockDanmaku(id string) (*model.Danmaku, error) {
	return s.transitionDanmaku(id, model.DanmakuStatusBlocked)
}

// BatchApproveDanmakus 批量审核通过弹幕，返回成功与失败数量。
func (s *Service) BatchApproveDanmakus(ids []string) (int, int, error) {
	success, failed := 0, 0
	for _, id := range ids {
		if _, err := s.transitionDanmaku(id, model.DanmakuStatusApproved); err != nil {
			failed++
			continue
		}
		success++
	}
	return success, failed, nil
}
