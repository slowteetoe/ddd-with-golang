package transport

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"slowteetoe.com/recommentations/recommendation/internal/recommendation"
)

func NewMux(recHandler recommendation.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/recommendation", recHandler.GetRecommendation)

	return r
}
