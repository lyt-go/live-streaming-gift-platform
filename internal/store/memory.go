package store

import (
	"sync"

	"livestream/internal/model"
)

// MemoryStore 是基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu          sync.RWMutex
	users       map[string]*model.User
	rooms       map[string]*model.Room
	danmakus    map[string]*model.Danmaku
	gifts       map[string]*model.Gift
	follows     map[string]*model.Follow
	giftRecords map[string]*model.GiftRecord
}

// NewMemoryStore 创建一个空的 MemoryStore。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:       make(map[string]*model.User),
		rooms:       make(map[string]*model.Room),
		danmakus:    make(map[string]*model.Danmaku),
		gifts:       make(map[string]*model.Gift),
		follows:     make(map[string]*model.Follow),
		giftRecords: make(map[string]*model.GiftRecord),
	}
}

var _ Store = (*MemoryStore)(nil)
