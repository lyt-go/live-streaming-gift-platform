package handler

import (
	"net/http"

	"livestream/internal/model"
	"livestream/pkg/httpx"
)

func (s *Server) registerDanmakuRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/danmakus", s.createDanmaku)
	mux.HandleFunc("GET /api/danmakus", s.listDanmakus)
	mux.HandleFunc("GET /api/danmakus/{id}", s.getDanmaku)
	mux.HandleFunc("POST /api/danmakus/{id}/approve", s.approveDanmaku)
	mux.HandleFunc("POST /api/danmakus/{id}/block", s.blockDanmaku)
	mux.HandleFunc("POST /api/danmakus/batch-approve", s.batchApproveDanmakus)
}

type createDanmakuRequest struct {
	RoomID  string `json:"room_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

func (s *Server) createDanmaku(w http.ResponseWriter, r *http.Request) {
	var req createDanmakuRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDanmaku(model.Danmaku{
		RoomID:  req.RoomID,
		UserID:  req.UserID,
		Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDanmakus(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DanmakuFilter{
		RoomID:  r.URL.Query().Get("room_id"),
		UserID:  r.URL.Query().Get("user_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDanmakus(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDanmaku(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDanmaku(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) approveDanmaku(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.ApproveDanmaku(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) blockDanmaku(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.BlockDanmaku(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type batchApproveRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchApproveDanmakus(w http.ResponseWriter, r *http.Request) {
	var req batchApproveRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	success, failed, err := s.svc.BatchApproveDanmakus(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"success": success, "failed": failed})
}
