package service

import (
	"livestream/internal/config"
	"livestream/internal/store"
	"livestream/pkg/logger"
)

// Service 聚合全部业务逻辑，依赖 Store 接口。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 构造 Service。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
