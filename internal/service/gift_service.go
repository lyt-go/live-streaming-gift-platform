package service

import (
	"sort"
	"time"

	"livestream/internal/model"
	"livestream/pkg/idgen"
)

// CreateGift 创建礼物。
func (s *Service) CreateGift(input model.Gift) (*model.Gift, error) {
	g := &model.Gift{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Icon:      input.Icon,
		Price:     input.Price,
		Category:  input.Category,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := g.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateGift(g); err != nil {
		return nil, err
	}
	return g, nil
}

// GetGift 按 ID 查询礼物。
func (s *Service) GetGift(id string) (*model.Gift, error) {
	return s.store.GetGift(id)
}

// ListGifts 分页查询礼物列表。
func (s *Service) ListGifts(filter model.GiftFilter, page, size int) ([]*model.Gift, int, error) {
	all := s.store.ListGifts()
	matched := make([]*model.Gift, 0, len(all))
	for _, g := range all {
		if filter.Match(g) {
			matched = append(matched, g)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Price > matched[j].Price
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Gift{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateGift 更新礼物信息。
func (s *Service) UpdateGift(id string, input model.Gift) (*model.Gift, error) {
	existing, err := s.store.GetGift(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Icon = input.Icon
	existing.Price = input.Price
	existing.Category = input.Category
	if input.Status != "" {
		existing.Status = input.Status
	}
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateGift(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteGift 删除礼物。
func (s *Service) DeleteGift(id string) error {
	return s.store.DeleteGift(id)
}

// BatchSetGiftStatus 批量上下架礼物。
func (s *Service) BatchSetGiftStatus(ids []string, status string) (int, int, error) {
	if status != model.GiftStatusActive && status != model.GiftStatusInactive {
		return 0, len(ids), model.NewValidationError("status", "礼物状态不合法")
	}
	success, failed := 0, 0
	for _, id := range ids {
		g, err := s.store.GetGift(id)
		if err != nil {
			failed++
			continue
		}
		g.Status = status
		g.UpdatedAt = time.Now()
		if err := s.store.UpdateGift(g); err != nil {
			failed++
			continue
		}
		success++
	}
	return success, failed, nil
}
