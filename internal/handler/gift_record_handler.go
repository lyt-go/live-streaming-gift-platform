package handler

import (
	"net/http"

	"livestream/internal/model"
	"livestream/pkg/httpx"
)

func (s *Server) registerGiftRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/gift-records", s.createGiftRecord)
	mux.HandleFunc("GET /api/gift-records", s.listGiftRecords)
	mux.HandleFunc("GET /api/gift-records/{id}", s.getGiftRecord)
}

type createGiftRecordRequest struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	GiftID   string `json:"gift_id"`
	Quantity int    `json:"quantity"`
}

func (s *Server) createGiftRecord(w http.ResponseWriter, r *http.Request) {
	var req createGiftRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateGiftRecord(model.GiftRecord{
		RoomID:   req.RoomID,
		UserID:   req.UserID,
		GiftID:   req.GiftID,
		Quantity: req.Quantity,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listGiftRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.GiftRecordFilter{
		RoomID: r.URL.Query().Get("room_id"),
		UserID: r.URL.Query().Get("user_id"),
		GiftID: r.URL.Query().Get("gift_id"),
	}
	items, total, err := s.svc.ListGiftRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getGiftRecord(w http.ResponseWriter, r *http.Request) {
	rec, err := s.svc.GetGiftRecord(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}
