package service

import (
	"sort"
	"strings"
	"time"

	"livestream/internal/model"
	"livestream/pkg/idgen"
)

// CreateUser 创建用户。
func (s *Service) CreateUser(input model.User) (*model.User, error) {
	u := &model.User{
		ID:             idgen.Hex(),
		Nickname:       input.Nickname,
		Avatar:         input.Avatar,
		Role:           input.Role,
		Bio:            input.Bio,
		FollowersCount: 0,
		Status:         input.Status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetUser 按 ID 查询用户。
func (s *Service) GetUser(id string) (*model.User, error) {
	return s.store.GetUser(id)
}

// ListUsers 分页查询用户列表。
func (s *Service) ListUsers(filter model.UserFilter, page, size int) ([]*model.User, int, error) {
	all := s.store.ListUsers()
	matched := make([]*model.User, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].FollowersCount != matched[j].FollowersCount {
			return matched[i].FollowersCount > matched[j].FollowersCount
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.User{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateUser 更新用户资料。
func (s *Service) UpdateUser(id string, input model.User) (*model.User, error) {
	existing, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	existing.Nickname = input.Nickname
	existing.Avatar = input.Avatar
	existing.Bio = input.Bio
	if strings.TrimSpace(input.Role) != "" {
		existing.Role = input.Role
	}
	if strings.TrimSpace(input.Status) != "" {
		existing.Status = input.Status
	}
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateUser(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteUser 删除用户，并清理其关联的全部关注关系。
// 关注关系涉及两端，被删用户既可能是关注者（A 关注 B）也可能是被关注者（C 关注 A），
// 因此删除时需要双向清理，并同步修正相关用户的粉丝计数，避免关注列表残留悬空关系。
func (s *Service) DeleteUser(id string) error {
	if _, err := s.store.GetUser(id); err != nil {
		return err
	}
	// 删除该用户相关的全部关注关系（作为关注者或被关注者）。
	removed := s.store.DeleteFollowsByUser(id)
	// 按 followeeID 汇总粉丝计数应减少的数量：
	// 仅在被删用户作为关注者（A 关注 B）时，B 的粉丝数需要减少；
	// 当用户作为被关注者被删除（C 关注 A）时，A 已删除，无需更新其计数。
	decrements := make(map[string]int64)
	for _, f := range removed {
		if f.FollowerID == id {
			decrements[f.FolloweeID]++
		}
	}
	for followeeID, n := range decrements {
		if followee, err := s.store.GetUser(followeeID); err == nil {
			if followee.FollowersCount > n {
				followee.FollowersCount -= n
			} else {
				followee.FollowersCount = 0
			}
			followee.UpdatedAt = time.Now()
			_ = s.store.UpdateUser(followee)
		}
	}
	return s.store.DeleteUser(id)
}

// BanUser 封禁用户。
func (s *Service) BanUser(id string) (*model.User, error) {
	u, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	if u.Status == model.UserStatusBanned {
		return nil, model.NewValidationError("status", "用户已处于封禁状态")
	}
	u.Status = model.UserStatusBanned
	u.UpdatedAt = time.Now()
	if err := s.store.UpdateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}
