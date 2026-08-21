package handler

import (
	"net/http"
	"strconv"

	"livestream/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/top-streamers", s.statsTopStreamers)
	mux.HandleFunc("GET /api/stats/room-danmakus", s.statsRoomDanmakus)
	mux.HandleFunc("GET /api/stats/gift-income", s.statsGiftIncome)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.Overview())
}

func (s *Server) statsTopStreamers(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	httpx.OK(w, s.svc.TopStreamers(limit))
}

func (s *Server) statsRoomDanmakus(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RoomDanmakuStats())
}

func (s *Server) statsGiftIncome(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.GiftIncomeByDay())
}
