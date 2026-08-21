package handler

import (
	"net/http"

	"livestream/internal/model"
	"livestream/pkg/httpx"
)

func (s *Server) registerRoomRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rooms", s.createRoom)
	mux.HandleFunc("GET /api/rooms", s.listRooms)
	mux.HandleFunc("GET /api/rooms/{id}", s.getRoom)
	mux.HandleFunc("PUT /api/rooms/{id}", s.updateRoom)
	mux.HandleFunc("DELETE /api/rooms/{id}", s.deleteRoom)
	mux.HandleFunc("POST /api/rooms/{id}/start", s.startRoom)
	mux.HandleFunc("POST /api/rooms/{id}/end", s.endRoom)
}

type createRoomRequest struct {
	Title      string `json:"title"`
	Category   string `json:"category"`
	StreamerID string `json:"streamer_id"`
	CoverURL   string `json:"cover_url"`
}

func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	room, err := s.svc.CreateRoom(model.Room{
		Title:      req.Title,
		Category:   req.Category,
		StreamerID: req.StreamerID,
		CoverURL:   req.CoverURL,
		Status:     model.RoomStatusPending,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, room)
}

func (s *Server) listRooms(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RoomFilter{
		Category:   r.URL.Query().Get("category"),
		Status:     r.URL.Query().Get("status"),
		StreamerID: r.URL.Query().Get("streamer_id"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRooms(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRoom(w http.ResponseWriter, r *http.Request) {
	room, err := s.svc.GetRoom(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, room)
}

func (s *Server) updateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	room, err := s.svc.UpdateRoom(r.PathValue("id"), model.Room{
		Title:    req.Title,
		Category: req.Category,
		CoverURL: req.CoverURL,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, room)
}

func (s *Server) deleteRoom(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRoom(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) startRoom(w http.ResponseWriter, r *http.Request) {
	room, err := s.svc.StartRoom(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, room)
}

func (s *Server) endRoom(w http.ResponseWriter, r *http.Request) {
	room, err := s.svc.EndRoom(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, room)
}
