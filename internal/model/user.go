package model

import (
	"strings"
	"time"
)

const (
	UserRoleStreamer = "streamer"
	UserRoleViewer   = "viewer"

	UserStatusActive = "active"
	UserStatusBanned = "banned"
)

// User 表示直播平台的用户，既可以是主播也可以是观众。
type User struct {
	ID             string    `json:"id"`
	Nickname       string    `json:"nickname"`
	Avatar         string    `json:"avatar"`
	Role           string    `json:"role"`
	Bio            string    `json:"bio"`
	FollowersCount int64     `json:"followers_count"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Role = strings.TrimSpace(u.Role)
	if u.Nickname == "" {
		return NewValidationError("nickname", "昵称不能为空")
	}
	if len([]rune(u.Nickname)) > 32 {
		return NewValidationError("nickname", "昵称长度不能超过 32 个字符")
	}
	if u.Role == "" {
		u.Role = UserRoleViewer
	}
	if u.Role != UserRoleStreamer && u.Role != UserRoleViewer {
		return NewValidationError("role", "用户角色不合法")
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}
	if u.Status != UserStatusActive && u.Status != UserStatusBanned {
		return NewValidationError("status", "用户状态不合法")
	}
	if u.FollowersCount < 0 {
		return NewValidationError("followers_count", "粉丝数不能为负数")
	}
	return nil
}

// UserFilter 用户列表筛选条件。
type UserFilter struct {
	Role    string
	Status  string
	Keyword string
}

func (f UserFilter) Match(u *User) bool {
	if f.Role != "" && u.Role != f.Role {
		return false
	}
	if f.Status != "" && u.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(u.Nickname), k) &&
			!strings.Contains(strings.ToLower(u.Bio), k) {
			return false
		}
	}
	return true
}
