package service

import (
	"sort"
	"time"

	"livestream/internal/model"
	"livestream/internal/store"
	"livestream/pkg/idgen"
)

// CreateGiftRecord 送礼，做跨实体校验：房间直播中、礼物有效、用户存在。
func (s *Service) CreateGiftRecord(input model.GiftRecord) (*model.GiftRecord, error) {
	room, err := s.store.GetRoom(input.RoomID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("room_id", "直播间不存在")
		}
		return nil, err
	}
	if room.Status != model.RoomStatusLive {
		return nil, model.NewValidationError("room_id", "直播间未在直播中")
	}
	gift, err := s.store.GetGift(input.GiftID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("gift_id", "礼物不存在")
		}
		return nil, err
	}
	if gift.Status != model.GiftStatusActive {
		return nil, model.NewValidationError("gift_id", "礼物已下架")
	}
	if _, err := s.store.GetUser(input.UserID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("user_id", "送礼用户不存在")
		}
		return nil, err
	}
	if input.Quantity <= 0 {
		return nil, model.NewValidationError("quantity", "赠送数量必须大于 0")
	}
	rec := &model.GiftRecord{
		ID:        idgen.Hex(),
		RoomID:    input.RoomID,
		UserID:    input.UserID,
		GiftID:    input.GiftID,
		Quantity:  input.Quantity,
		Amount:    gift.Price * int64(input.Quantity),
		CreatedAt: time.Now(),
	}
	if err := rec.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateGiftRecord(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// GetGiftRecord 按 ID 查询送礼记录。
func (s *Service) GetGiftRecord(id string) (*model.GiftRecord, error) {
	return s.store.GetGiftRecord(id)
}

// ListGiftRecords 分页查询送礼记录。
func (s *Service) ListGiftRecords(filter model.GiftRecordFilter, page, size int) ([]*model.GiftRecord, int, error) {
	all := s.store.ListGiftRecords()
	matched := make([]*model.GiftRecord, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.GiftRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
