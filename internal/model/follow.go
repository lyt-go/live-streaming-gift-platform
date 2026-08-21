package model

import (
	"strings"
	"time"
)

// Follow 表示一条关注关系，follower 关注 followee。
type Follow struct {
	ID         string    `json:"id"`
	FollowerID string    `json:"follower_id"`
	FolloweeID string    `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (f *Follow) Validate() error {
	f.FollowerID = strings.TrimSpace(f.FollowerID)
	f.FolloweeID = strings.TrimSpace(f.FolloweeID)
	if f.FollowerID == "" {
		return NewValidationError("follower_id", "关注者不能为空")
	}
	if f.FolloweeID == "" {
		return NewValidationError("followee_id", "被关注者不能为空")
	}
	if f.FollowerID == f.FolloweeID {
		return NewValidationError("followee_id", "不能关注自己")
	}
	return nil
}

// FollowFilter 关注关系筛选条件。
type FollowFilter struct {
	FollowerID string
	FolloweeID string
}

func (f FollowFilter) Match(v *Follow) bool {
	if f.FollowerID != "" && v.FollowerID != f.FollowerID {
		return false
	}
	if f.FolloweeID != "" && v.FolloweeID != f.FolloweeID {
		return false
	}
	return true
}
