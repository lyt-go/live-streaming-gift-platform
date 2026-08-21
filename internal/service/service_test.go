package service

import (
	"testing"

	"livestream/internal/config"
	"livestream/internal/model"
	"livestream/internal/store"
	"livestream/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func mustCreateUser(t *testing.T, s *Service, nickname, role string) *model.User {
	t.Helper()
	u, err := s.CreateUser(model.User{Nickname: nickname, Role: role})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	return u
}

func mustCreateGift(t *testing.T, s *Service, name string, price int64) *model.Gift {
	t.Helper()
	g, err := s.CreateGift(model.Gift{Name: name, Price: price})
	if err != nil {
		t.Fatalf("CreateGift: %v", err)
	}
	return g
}

func TestRoomLifecycle(t *testing.T) {
	s := newTestService()
	streamer := mustCreateUser(t, s, "主播", model.UserRoleStreamer)
	room, err := s.CreateRoom(model.Room{Title: "直播", Category: "游戏", StreamerID: streamer.ID})
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	// 开播
	room, err = s.StartRoom(room.ID)
	if err != nil || room.Status != model.RoomStatusLive {
		t.Fatalf("StartRoom: %v %v", room, err)
	}
	// 已直播不能再开播
	if _, err := s.StartRoom(room.ID); err == nil {
		t.Fatalf("期望开播失败")
	}
	// 下播
	room, err = s.EndRoom(room.ID)
	if err != nil || room.Status != model.RoomStatusEnded {
		t.Fatalf("EndRoom: %v %v", room, err)
	}
}

func TestDanmakuModeration(t *testing.T) {
	s := newTestService()
	streamer := mustCreateUser(t, s, "主播", model.UserRoleStreamer)
	viewer := mustCreateUser(t, s, "观众", model.UserRoleViewer)
	room, _ := s.CreateRoom(model.Room{Title: "直播", Category: "游戏", StreamerID: streamer.ID})
	d, err := s.CreateDanmaku(model.Danmaku{RoomID: room.ID, UserID: viewer.ID, Content: "哈哈"})
	if err != nil {
		t.Fatalf("CreateDanmaku: %v", err)
	}
	if d.Status != model.DanmakuStatusPending {
		t.Fatalf("期望 pending, got %s", d.Status)
	}
	if _, err := s.ApproveDanmaku(d.ID); err != nil {
		t.Fatalf("ApproveDanmaku: %v", err)
	}
	// 已通过不能再拦截
	if _, err := s.BlockDanmaku(d.ID); err == nil {
		t.Fatalf("期望拦截失败")
	}
}

func TestDeletingDanmakuSenderRemovesDanmakus(t *testing.T) {
	s := newTestService()
	user := mustCreateUser(t, s, "danmaku-delete", model.UserRoleViewer)
	streamer := mustCreateUser(t, s, "danmaku-streamer", model.UserRoleStreamer)
	room, err := s.CreateRoom(model.Room{Title: "danmaku-room", Category: "game", StreamerID: streamer.ID})
	if err != nil { t.Fatal(err) }
	if _, err = s.CreateDanmaku(model.Danmaku{RoomID: room.ID, UserID: user.ID, Content: "orphan"}); err != nil { t.Fatal(err) }
	if err = s.DeleteUser(user.ID); err != nil { t.Fatal(err) }
	items, total, err := s.ListDanmakus(model.DanmakuFilter{UserID: user.ID}, 1, 10)
	if err != nil { t.Fatal(err) }
	if total != 0 || len(items) != 0 { t.Fatalf("expected danmakus removed, got total=%d len=%d", total, len(items)) }
}

func TestGiftRecordCrossEntityValidation(t *testing.T) {
	s := newTestService()
	streamer := mustCreateUser(t, s, "主播", model.UserRoleStreamer)
	viewer := mustCreateUser(t, s, "观众", model.UserRoleViewer)
	gift := mustCreateGift(t, s, "火箭", 10000)
	room, _ := s.CreateRoom(model.Room{Title: "直播", Category: "游戏", StreamerID: streamer.ID})

	// 未开播不能送礼
	if _, err := s.CreateGiftRecord(model.GiftRecord{RoomID: room.ID, UserID: viewer.ID, GiftID: gift.ID, Quantity: 1}); err == nil {
		t.Fatalf("期望未开播送礼失败")
	}
	// 开播后送礼
	room, _ = s.StartRoom(room.ID)
	rec, err := s.CreateGiftRecord(model.GiftRecord{RoomID: room.ID, UserID: viewer.ID, GiftID: gift.ID, Quantity: 3})
	if err != nil {
		t.Fatalf("CreateGiftRecord: %v", err)
	}
	if rec.Amount != 30000 {
		t.Fatalf("期望金额 30000, got %d", rec.Amount)
	}
}

func TestFollowIncrementsCount(t *testing.T) {
	s := newTestService()
	a := mustCreateUser(t, s, "a", model.UserRoleViewer)
	b := mustCreateUser(t, s, "b", model.UserRoleStreamer)
	if _, err := s.Follow(a.ID, b.ID); err != nil {
		t.Fatalf("Follow: %v", err)
	}
	got, _ := s.GetUser(b.ID)
	if got.FollowersCount != 1 {
		t.Fatalf("期望粉丝数 1, got %d", got.FollowersCount)
	}
	// 重复关注冲突
	if _, err := s.Follow(a.ID, b.ID); err == nil {
		t.Fatalf("期望重复关注失败")
	}
	if err := s.Unfollow(a.ID, b.ID); err != nil {
		t.Fatalf("Unfollow: %v", err)
	}
	got, _ = s.GetUser(b.ID)
	if got.FollowersCount != 0 {
		t.Fatalf("期望粉丝数 0, got %d", got.FollowersCount)
	}
}

func TestTopStreamers(t *testing.T) {
	s := newTestService()
	s1 := mustCreateUser(t, s, "主播1", model.UserRoleStreamer)
	s2 := mustCreateUser(t, s, "主播2", model.UserRoleStreamer)
	viewer := mustCreateUser(t, s, "观众", model.UserRoleViewer)
	gift := mustCreateGift(t, s, "火箭", 10000)

	room1, _ := s.CreateRoom(model.Room{Title: "r1", Category: "游戏", StreamerID: s1.ID})
	room2, _ := s.CreateRoom(model.Room{Title: "r2", Category: "游戏", StreamerID: s2.ID})
	room1, _ = s.StartRoom(room1.ID)
	room2, _ = s.StartRoom(room2.ID)

	_, _ = s.CreateGiftRecord(model.GiftRecord{RoomID: room1.ID, UserID: viewer.ID, GiftID: gift.ID, Quantity: 1})
	_, _ = s.CreateGiftRecord(model.GiftRecord{RoomID: room2.ID, UserID: viewer.ID, GiftID: gift.ID, Quantity: 5})

	ranks := s.TopStreamers(10)
	if len(ranks) != 2 {
		t.Fatalf("期望 2 个主播, got %d", len(ranks))
	}
	if ranks[0].StreamerID != s2.ID || ranks[0].Amount != 50000 {
		t.Fatalf("期望主播2 排第一, got %+v", ranks[0])
	}
}
