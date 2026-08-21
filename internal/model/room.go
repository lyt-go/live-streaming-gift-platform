package model

import (
	"strings"
	"time"
)

const (
	RoomStatusPending = "pending"
	RoomStatusLive    = "live"
	RoomStatusEnded   = "ended"
)

// roomTransitions 定义直播间的合法状态流转。
var roomTransitions = map[string]map[string]bool{
	RoomStatusPending: {RoomStatusLive: true},
	RoomStatusLive:    {RoomStatusEnded: true},
	RoomStatusEnded:   {},
}

// CanTransition 判断直播间能否从 from 状态流转到 to 状态。
func CanTransition(from, to string) bool {
	if m, ok := roomTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Room 表示一个直播间。
type Room struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	StreamerID  string    `json:"streamer_id"`
	CoverURL    string    `json:"cover_url"`
	Status      string    `json:"status"`
	ViewerCount int64     `json:"viewer_count"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Room) Validate() error {
	r.Title = strings.TrimSpace(r.Title)
	r.Category = strings.TrimSpace(r.Category)
	r.StreamerID = strings.TrimSpace(r.StreamerID)
	if r.Title == "" {
		return NewValidationError("title", "直播间标题不能为空")
	}
	if len([]rune(r.Title)) > 100 {
		return NewValidationError("title", "直播间标题长度不能超过 100 个字符")
	}
	if r.Category == "" {
		return NewValidationError("category", "直播分类不能为空")
	}
	if r.StreamerID == "" {
		return NewValidationError("streamer_id", "主播不能为空")
	}
	if r.Status == "" {
		r.Status = RoomStatusPending
	}
	if r.Status != RoomStatusPending && r.Status != RoomStatusLive && r.Status != RoomStatusEnded {
		return NewValidationError("status", "直播间状态不合法")
	}
	if r.ViewerCount < 0 {
		return NewValidationError("viewer_count", "观看人数不能为负数")
	}
	return nil
}

// RoomFilter 直播间列表筛选条件。
type RoomFilter struct {
	Category   string
	Status     string
	StreamerID string
	Keyword    string
}

func (f RoomFilter) Match(r *Room) bool {
	if f.Category != "" && r.Category != f.Category {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.StreamerID != "" && r.StreamerID != f.StreamerID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Title), k) {
			return false
		}
	}
	return true
}
