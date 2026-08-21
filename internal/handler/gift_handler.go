package handler

import (
	"net/http"

	"livestream/internal/model"
	"livestream/pkg/httpx"
)

func (s *Server) registerGiftRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/gifts", s.createGift)
	mux.HandleFunc("GET /api/gifts", s.listGifts)
	mux.HandleFunc("GET /api/gifts/{id}", s.getGift)
	mux.HandleFunc("PUT /api/gifts/{id}", s.updateGift)
	mux.HandleFunc("DELETE /api/gifts/{id}", s.deleteGift)
	mux.HandleFunc("POST /api/gifts/batch-status", s.batchSetGiftStatus)
}

type createGiftRequest struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Price    int64  `json:"price"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

func (s *Server) createGift(w http.ResponseWriter, r *http.Request) {
	var req createGiftRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	g, err := s.svc.CreateGift(model.Gift{
		Name:     req.Name,
		Icon:     req.Icon,
		Price:    req.Price,
		Category: req.Category,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, g)
}

func (s *Server) listGifts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.GiftFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListGifts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getGift(w http.ResponseWriter, r *http.Request) {
	g, err := s.svc.GetGift(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, g)
}

func (s *Server) updateGift(w http.ResponseWriter, r *http.Request) {
	var req createGiftRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	g, err := s.svc.UpdateGift(r.PathValue("id"), model.Gift{
		Name:     req.Name,
		Icon:     req.Icon,
		Price:    req.Price,
		Category: req.Category,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, g)
}

func (s *Server) deleteGift(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteGift(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchGiftStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

func (s *Server) batchSetGiftStatus(w http.ResponseWriter, r *http.Request) {
	var req batchGiftStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	success, failed, err := s.svc.BatchSetGiftStatus(req.IDs, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"success": success, "failed": failed})
}
