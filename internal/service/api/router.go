package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

// IService describes all the external methods that can be used in the API handlers.
type IService interface{}

func router(log *logan.Entry, svc IService) chi.Router {
	r := chi.NewRouter()

	r.Use(
		middleware.Heartbeat("/health"),
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			CtxLog(log),
			CtxService(svc),
		),
	)

	r.Route("/api/v1", func(r chi.Router) {})

	return r
}
