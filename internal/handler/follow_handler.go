package handler

import (
	"net/http"

	"livestream/internal/model"
	"livestream/pkg/httpx"
)

func (s *Server) registerFollowRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/follows", s.follow)
	mux.HandleFunc("DELETE /api/follows", s.unfollow)
	mux.HandleFunc("GET /api/follows", s.listFollows)
	mux.HandleFunc("GET /api/follows/{id}", s.getFollow)
}

type followRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

func (s *Server) follow(w http.ResponseWriter, r *http.Request) {
	var req followRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.Follow(req.FollowerID, req.FolloweeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) unfollow(w http.ResponseWriter, r *http.Request) {
	var req followRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.Unfollow(req.FollowerID, req.FolloweeID); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listFollows(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FollowFilter{
		FollowerID: r.URL.Query().Get("follower_id"),
		FolloweeID: r.URL.Query().Get("followee_id"),
	}
	items, total, err := s.svc.ListFollows(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFollow(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.GetFollow(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}
