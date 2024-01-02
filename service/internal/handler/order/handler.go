package order

import "github.com/Nahbox/streamed-order-viewer/service/internal/cached"

type Handler struct {
	DB           DB
	CacheManager *cached.CacheManager
}

func NewHandler(db DB, cacheManager *cached.CacheManager) *Handler {
	handler := &Handler{
		DB:           db,
		CacheManager: cacheManager,
	}

	return handler
}
