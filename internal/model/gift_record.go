package model

import (
	"strings"
	"time"
)

// GiftRecord 表示一条送礼记录，amount 为金额，单位分。
type GiftRecord struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	GiftID    string    `json:"gift_id"`
	Quantity  int       `json:"quantity"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *GiftRecord) Validate() error {
	r.RoomID = strings.TrimSpace(r.RoomID)
	r.UserID = strings.TrimSpace(r.UserID)
	r.GiftID = strings.TrimSpace(r.GiftID)
	if r.RoomID == "" {
		return NewValidationError("room_id", "直播间不能为空")
	}
	if r.UserID == "" {
		return NewValidationError("user_id", "送礼用户不能为空")
	}
	if r.GiftID == "" {
		return NewValidationError("gift_id", "礼物不能为空")
	}
	if r.Quantity <= 0 {
		return NewValidationError("quantity", "赠送数量必须大于 0")
	}
	if r.Amount <= 0 {
		return NewValidationError("amount", "金额必须大于 0")
	}
	return nil
}

// GiftRecordFilter 送礼记录筛选条件。
type GiftRecordFilter struct {
	RoomID string
	UserID string
	GiftID string
}

func (f GiftRecordFilter) Match(r *GiftRecord) bool {
	if f.RoomID != "" && r.RoomID != f.RoomID {
		return false
	}
	if f.UserID != "" && r.UserID != f.UserID {
		return false
	}
	if f.GiftID != "" && r.GiftID != f.GiftID {
		return false
	}
	return true
}
