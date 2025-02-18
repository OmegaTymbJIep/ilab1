package api

import (
	"github.com/go-chi/chi/v5"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

// IService describes all the external methods that can be used in the API handlers.
type IService interface{}

func initRouter(log *logan.Entry, svc IService) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			CtxLog(log),
			CtxService(svc),
		),
	)

	return r
}

func (a *api) finalRouter() chi.Router {
	r := a.router.Route("/v1", func(r chi.Router) {})

	a.router.Mount("/api", r)

	return a.router
}
