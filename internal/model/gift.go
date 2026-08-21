package model

import (
	"strings"
	"time"
)

const (
	GiftStatusActive   = "active"
	GiftStatusInactive = "inactive"
)

// Gift 表示一种可赠送的虚拟礼物，金额单位为分。
type Gift struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Price     int64     `json:"price"`
	Category  string    `json:"category"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *Gift) Validate() error {
	g.Name = strings.TrimSpace(g.Name)
	g.Category = strings.TrimSpace(g.Category)
	if g.Name == "" {
		return NewValidationError("name", "礼物名称不能为空")
	}
	if g.Price <= 0 {
		return NewValidationError("price", "礼物价格必须大于 0")
	}
	if g.Status == "" {
		g.Status = GiftStatusActive
	}
	if g.Status != GiftStatusActive && g.Status != GiftStatusInactive {
		return NewValidationError("status", "礼物状态不合法")
	}
	return nil
}

// GiftFilter 礼物列表筛选条件。
type GiftFilter struct {
	Category string
	Status   string
	Keyword  string
}

func (f GiftFilter) Match(g *Gift) bool {
	if f.Category != "" && g.Category != f.Category {
		return false
	}
	if f.Status != "" && g.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(g.Name), k) {
			return false
		}
	}
	return true
}
