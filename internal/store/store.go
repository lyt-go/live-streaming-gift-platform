// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"livestream/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// User 用户
	CreateUser(u *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByNickname(nickname string) (*model.User, error)
	ListUsers() []*model.User
	UpdateUser(u *model.User) error
	DeleteUser(id string) error

	// Room 直播间
	CreateRoom(r *model.Room) error
	GetRoom(id string) (*model.Room, error)
	ListRooms() []*model.Room
	UpdateRoom(r *model.Room) error
	DeleteRoom(id string) error

	// Danmaku 弹幕
	CreateDanmaku(d *model.Danmaku) error
	GetDanmaku(id string) (*model.Danmaku, error)
	ListDanmakus() []*model.Danmaku
	UpdateDanmaku(d *model.Danmaku) error
	DeleteDanmaku(id string) error

	// Gift 礼物
	CreateGift(g *model.Gift) error
	GetGift(id string) (*model.Gift, error)
	ListGifts() []*model.Gift
	UpdateGift(g *model.Gift) error
	DeleteGift(id string) error

	// Follow 关注关系
	CreateFollow(f *model.Follow) error
	GetFollow(id string) (*model.Follow, error)
	ListFollows() []*model.Follow
	DeleteFollow(id string) error
	FollowExists(followerID, followeeID string) bool

	// GiftRecord 送礼记录
	CreateGiftRecord(r *model.GiftRecord) error
	GetGiftRecord(id string) (*model.GiftRecord, error)
	ListGiftRecords() []*model.GiftRecord
	DeleteGiftRecord(id string) error
}
