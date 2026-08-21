package service

import (
	"sort"

	"livestream/internal/model"
)

// Overview 平台总览统计。
type Overview struct {
	UserCount       int   `json:"user_count"`
	StreamerCount   int   `json:"streamer_count"`
	RoomCount       int   `json:"room_count"`
	LiveRoomCount   int   `json:"live_room_count"`
	GiftCount       int   `json:"gift_count"`
	FollowCount     int   `json:"follow_count"`
	GiftRecordCount int   `json:"gift_record_count"`
	TotalGiftAmount int64 `json:"total_gift_amount"`
}

// Overview 计算平台总览统计。
func (s *Service) Overview() *Overview {
	o := &Overview{}
	for _, u := range s.store.ListUsers() {
		o.UserCount++
		if u.Role == model.UserRoleStreamer {
			o.StreamerCount++
		}
	}
	for _, r := range s.store.ListRooms() {
		o.RoomCount++
		if r.Status == model.RoomStatusLive {
			o.LiveRoomCount++
		}
	}
	o.GiftCount = len(s.store.ListGifts())
	o.FollowCount = len(s.store.ListFollows())
	for _, r := range s.store.ListGiftRecords() {
		o.GiftRecordCount++
		o.TotalGiftAmount += r.Amount
	}
	return o
}

// StreamerRank 主播打赏榜条目。
type StreamerRank struct {
	StreamerID string `json:"streamer_id"`
	Nickname   string `json:"nickname"`
	Amount     int64  `json:"amount"`
	Records    int    `json:"records"`
}

// TopStreamers 返回按收到礼物金额排序的主播打赏榜前 N 名。
func (s *Service) TopStreamers(limit int) []*StreamerRank {
	if limit <= 0 {
		limit = 10
	}
	roomOwner := map[string]string{}
	for _, r := range s.store.ListRooms() {
		roomOwner[r.ID] = r.StreamerID
	}
	agg := map[string]*StreamerRank{}
	for _, rec := range s.store.ListGiftRecords() {
		streamerID := roomOwner[rec.RoomID]
		if streamerID == "" {
			continue
		}
		if _, ok := agg[streamerID]; !ok {
			agg[streamerID] = &StreamerRank{StreamerID: streamerID}
		}
		agg[streamerID].Amount += rec.Amount
		agg[streamerID].Records++
	}
	ranks := make([]*StreamerRank, 0, len(agg))
	for _, rank := range agg {
		if u, err := s.store.GetUser(rank.StreamerID); err == nil {
			rank.Nickname = u.Nickname
		}
		ranks = append(ranks, rank)
	}
	sort.Slice(ranks, func(i, j int) bool {
		if ranks[i].Amount != ranks[j].Amount {
			return ranks[i].Amount > ranks[j].Amount
		}
		return ranks[i].Records > ranks[j].Records
	})
	if len(ranks) > limit {
		ranks = ranks[:limit]
	}
	return ranks
}

// RoomDanmakuStat 单个直播间的弹幕统计。
type RoomDanmakuStat struct {
	RoomID   string `json:"room_id"`
	Title    string `json:"title"`
	Total    int    `json:"total"`
	Pending  int    `json:"pending"`
	Approved int    `json:"approved"`
	Blocked  int    `json:"blocked"`
}

// RoomDanmakuStats 统计每个直播间的弹幕数量分布。
func (s *Service) RoomDanmakuStats() []*RoomDanmakuStat {
	stats := map[string]*RoomDanmakuStat{}
	for _, d := range s.store.ListDanmakus() {
		st, ok := stats[d.RoomID]
		if !ok {
			st = &RoomDanmakuStat{RoomID: d.RoomID}
			stats[d.RoomID] = st
		}
		st.Total++
		switch d.Status {
		case model.DanmakuStatusPending:
			st.Pending++
		case model.DanmakuStatusApproved:
			st.Approved++
		case model.DanmakuStatusBlocked:
			st.Blocked++
		}
	}
	out := make([]*RoomDanmakuStat, 0, len(stats))
	for _, st := range stats {
		if r, err := s.store.GetRoom(st.RoomID); err == nil {
			st.Title = r.Title
		}
		out = append(out, st)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Total > out[j].Total
	})
	return out
}

// DailyGiftIncome 按天的礼物收入。
type DailyGiftIncome struct {
	Date   string `json:"date"`
	Amount int64  `json:"amount"`
	Records int   `json:"records"`
}

// GiftIncomeByDay 按天汇总礼物收入。
func (s *Service) GiftIncomeByDay() []*DailyGiftIncome {
	agg := map[string]*DailyGiftIncome{}
	for _, rec := range s.store.ListGiftRecords() {
		day := rec.CreatedAt.Format("2006-01-02")
		item, ok := agg[day]
		if !ok {
			item = &DailyGiftIncome{Date: day}
			agg[day] = item
		}
		item.Amount += rec.Amount
		item.Records++
	}
	out := make([]*DailyGiftIncome, 0, len(agg))
	for _, item := range agg {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Date < out[j].Date
	})
	return out
}
