package model

import (
	"strings"
	"time"
)

const (
	DanmakuStatusPending  = "pending"
	DanmakuStatusApproved = "approved"
	DanmakuStatusBlocked  = "blocked"
)

// danmakuTransitions 定义弹幕审核状态流转。
var danmakuTransitions = map[string]map[string]bool{
	DanmakuStatusPending: {DanmakuStatusApproved: true, DanmakuStatusBlocked: true},
	DanmakuStatusApproved: {},
	DanmakuStatusBlocked:  {},
}

// CanTransitionDanmaku 判断弹幕能否从 from 流转到 to。
func CanTransitionDanmaku(from, to string) bool {
	if m, ok := danmakuTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Danmaku 表示一条弹幕。
type Danmaku struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *Danmaku) Validate() error {
	d.RoomID = strings.TrimSpace(d.RoomID)
	d.UserID = strings.TrimSpace(d.UserID)
	d.Content = strings.TrimSpace(d.Content)
	if d.RoomID == "" {
		return NewValidationError("room_id", "直播间不能为空")
	}
	if d.UserID == "" {
		return NewValidationError("user_id", "发送用户不能为空")
	}
	if d.Content == "" {
		return NewValidationError("content", "弹幕内容不能为空")
	}
	if len([]rune(d.Content)) > 200 {
		return NewValidationError("content", "弹幕内容长度不能超过 200 个字符")
	}
	if d.Status == "" {
		d.Status = DanmakuStatusPending
	}
	if d.Status != DanmakuStatusPending && d.Status != DanmakuStatusApproved && d.Status != DanmakuStatusBlocked {
		return NewValidationError("status", "弹幕状态不合法")
	}
	return nil
}

// DanmakuFilter 弹幕列表筛选条件。
type DanmakuFilter struct {
	RoomID  string
	UserID  string
	Status  string
	Keyword string
}

func (f DanmakuFilter) Match(d *Danmaku) bool {
	if f.RoomID != "" && d.RoomID != f.RoomID {
		return false
	}
	if f.UserID != "" && d.UserID != f.UserID {
		return false
	}
	if f.Status != "" && d.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Content), k) {
			return false
		}
	}
	return true
}
