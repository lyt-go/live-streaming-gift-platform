package service

import (
	"sort"
	"time"

	"livestream/internal/model"
	"livestream/internal/store"
	"livestream/pkg/idgen"
)

// Follow 建立关注关系。
func (s *Service) Follow(followerID, followeeID string) (*model.Follow, error) {
	if _, err := s.store.GetUser(followerID); err != nil {
		return nil, model.NewValidationError("follower_id", "关注者不存在")
	}
	if _, err := s.store.GetUser(followeeID); err != nil {
		return nil, model.NewValidationError("followee_id", "被关注者不存在")
	}
	f := &model.Follow{
		ID:         idgen.Hex(),
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  time.Now(),
	}
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if followee, err := s.store.GetUser(followeeID); err == nil {
		followee.FollowersCount++
		followee.UpdatedAt = time.Now()
		_ = s.store.UpdateUser(followee)
	}
	if err := s.store.CreateFollow(f); err != nil {
		return nil, err
	}
	// 关注成功后累加被关注者粉丝数。
	if followee, err := s.store.GetUser(followeeID); err == nil {
		followee.FollowersCount++
		followee.UpdatedAt = time.Now()
		_ = s.store.UpdateUser(followee)
	}
	return f, nil
}

// Unfollow 取消关注。
func (s *Service) Unfollow(followerID, followeeID string) error {
	var target *model.Follow
	for _, f := range s.store.ListFollows() {
		if f.FollowerID == followerID && f.FolloweeID == followeeID {
			target = f
			break
		}
	}
	if target == nil {
		return store.ErrNotFound
	}
	if err := s.store.DeleteFollow(target.ID); err != nil {
		return err
	}
	if followee, err := s.store.GetUser(followeeID); err == nil && followee.FollowersCount > 0 {
		followee.FollowersCount--
		followee.UpdatedAt = time.Now()
		_ = s.store.UpdateUser(followee)
	}
	return nil
}

// GetFollow 按 ID 查询关注关系。
func (s *Service) GetFollow(id string) (*model.Follow, error) {
	return s.store.GetFollow(id)
}

// ListFollows 分页查询关注关系。
func (s *Service) ListFollows(filter model.FollowFilter, page, size int) ([]*model.Follow, int, error) {
	all := s.store.ListFollows()
	matched := make([]*model.Follow, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Follow{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
