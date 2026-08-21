package store

import (
	"testing"
	"time"

	"livestream/internal/model"
)

func newTestStore() *MemoryStore {
	return NewMemoryStore()
}

func newUser(id, nickname, role string) *model.User {
	return &model.User{
		ID:        id,
		Nickname:  nickname,
		Role:      role,
		Status:    model.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestUserCRUD(t *testing.T) {
	s := newTestStore()
	u := newUser("u1", "alice", model.UserRoleStreamer)
	if err := s.CreateUser(u); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	// 重复昵称冲突
	if err := s.CreateUser(newUser("u2", "alice", model.UserRoleViewer)); err != ErrConflict {
		t.Fatalf("期望冲突, got %v", err)
	}
	got, err := s.GetUser("u1")
	if err != nil || got.Nickname != "alice" {
		t.Fatalf("GetUser: %v %v", got, err)
	}
	if _, err := s.GetUser("missing"); err != ErrNotFound {
		t.Fatalf("期望 NotFound, got %v", err)
	}
	got.Role = model.UserRoleViewer
	if err := s.UpdateUser(got); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if len(s.ListUsers()) != 1 {
		t.Fatalf("期望 1 个用户, got %d", len(s.ListUsers()))
	}
	if err := s.DeleteUser("u1"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if err := s.DeleteUser("u1"); err != ErrNotFound {
		t.Fatalf("期望 NotFound, got %v", err)
	}
}

func TestRoomCRUD(t *testing.T) {
	s := newTestStore()
	r := &model.Room{ID: "r1", Title: "测试房间", Category: "游戏", StreamerID: "u1", Status: model.RoomStatusPending}
	if err := s.CreateRoom(r); err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	got, err := s.GetRoom("r1")
	if err != nil || got.Title != "测试房间" {
		t.Fatalf("GetRoom: %v %v", got, err)
	}
	got.Status = model.RoomStatusLive
	if err := s.UpdateRoom(got); err != nil {
		t.Fatalf("UpdateRoom: %v", err)
	}
	if len(s.ListRooms()) != 1 {
		t.Fatalf("期望 1 个房间, got %d", len(s.ListRooms()))
	}
	if err := s.DeleteRoom("r1"); err != nil {
		t.Fatalf("DeleteRoom: %v", err)
	}
}

func TestDanmakuCRUD(t *testing.T) {
	s := newTestStore()
	d := &model.Danmaku{ID: "d1", RoomID: "r1", UserID: "u1", Content: "666", Status: model.DanmakuStatusPending}
	if err := s.CreateDanmaku(d); err != nil {
		t.Fatalf("CreateDanmaku: %v", err)
	}
	got, err := s.GetDanmaku("d1")
	if err != nil || got.Content != "666" {
		t.Fatalf("GetDanmaku: %v %v", got, err)
	}
	got.Status = model.DanmakuStatusApproved
	if err := s.UpdateDanmaku(got); err != nil {
		t.Fatalf("UpdateDanmaku: %v", err)
	}
	if err := s.DeleteDanmaku("d1"); err != nil {
		t.Fatalf("DeleteDanmaku: %v", err)
	}
}

func TestGiftCRUD(t *testing.T) {
	s := newTestStore()
	g := &model.Gift{ID: "g1", Name: "火箭", Price: 10000, Status: model.GiftStatusActive}
	if err := s.CreateGift(g); err != nil {
		t.Fatalf("CreateGift: %v", err)
	}
	if err := s.CreateGift(&model.Gift{ID: "g2", Name: "火箭", Price: 10000}); err != ErrConflict {
		t.Fatalf("期望冲突, got %v", err)
	}
	if _, err := s.GetGift("g1"); err != nil {
		t.Fatalf("GetGift: %v", err)
	}
	if err := s.DeleteGift("g1"); err != nil {
		t.Fatalf("DeleteGift: %v", err)
	}
}

func TestFollowUniqueAndExists(t *testing.T) {
	s := newTestStore()
	f := &model.Follow{ID: "f1", FollowerID: "u1", FolloweeID: "u2"}
	if err := s.CreateFollow(f); err != nil {
		t.Fatalf("CreateFollow: %v", err)
	}
	if err := s.CreateFollow(&model.Follow{ID: "f2", FollowerID: "u1", FolloweeID: "u2"}); err != ErrConflict {
		t.Fatalf("期望冲突, got %v", err)
	}
	if !s.FollowExists("u1", "u2") {
		t.Fatalf("期望关注关系存在")
	}
	if s.FollowExists("u2", "u1") {
		t.Fatalf("期望关注关系不存在")
	}
}

func TestDeleteFollowsByUser(t *testing.T) {
	s := newTestStore()
	// u1 关注 u2、u1 关注 u3、u3 关注 u1：删除 u1 应同时清理三端关系。
	if err := s.CreateFollow(&model.Follow{ID: "f1", FollowerID: "u1", FolloweeID: "u2"}); err != nil {
		t.Fatalf("CreateFollow: %v", err)
	}
	if err := s.CreateFollow(&model.Follow{ID: "f2", FollowerID: "u1", FolloweeID: "u3"}); err != nil {
		t.Fatalf("CreateFollow: %v", err)
	}
	if err := s.CreateFollow(&model.Follow{ID: "f3", FollowerID: "u3", FolloweeID: "u1"}); err != nil {
		t.Fatalf("CreateFollow: %v", err)
	}
	removed := s.DeleteFollowsByUser("u1")
	if len(removed) != 3 {
		t.Fatalf("期望删除 3 条关系, got %d", len(removed))
	}
	for _, f := range removed {
		if f.FollowerID != "u1" && f.FolloweeID != "u1" {
			t.Fatalf("删除了无关关系: %+v", f)
		}
	}
	if len(s.ListFollows()) != 0 {
		t.Fatalf("期望关注关系全部清理, got %d", len(s.ListFollows()))
	}
	// 无关联用户删除应返回空切片。
	if r := s.DeleteFollowsByUser("no-such-user"); len(r) != 0 {
		t.Fatalf("期望空切片, got %d", len(r))
	}
}

func TestGiftRecordCRUD(t *testing.T) {
	s := newTestStore()
	rec := &model.GiftRecord{ID: "gr1", RoomID: "r1", UserID: "u1", GiftID: "g1", Quantity: 2, Amount: 20000}
	if err := s.CreateGiftRecord(rec); err != nil {
		t.Fatalf("CreateGiftRecord: %v", err)
	}
	if _, err := s.GetGiftRecord("gr1"); err != nil {
		t.Fatalf("GetGiftRecord: %v", err)
	}
	if err := s.DeleteGiftRecord("gr1"); err != nil {
		t.Fatalf("DeleteGiftRecord: %v", err)
	}
}
