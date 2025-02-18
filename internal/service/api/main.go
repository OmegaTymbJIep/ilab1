package api

import (
	"context"
	"net"

	"github.com/go-chi/chi/v5"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

type API interface {
	Run(context.Context)
	Router() chi.Router
}

type api struct {
	log      *logan.Entry
	listener net.Listener

	router chi.Router
}

func New(log *logan.Entry, listener net.Listener, service IService) (API, error) {
	return &api{
		log:      log,
		listener: listener,

		router: initRouter(log, service),
	}, nil
}

func (a *api) Run(ctx context.Context) {
	a.log.Info("API started")
	ape.Serve(ctx, a.finalRouter(), a, ape.ServeOpts{})
}

func (a *api) Log() *logan.Entry {
	return a.log
}

func (a *api) Listener() net.Listener {
	return a.listener
}

func (a *api) Router() chi.Router {
	return a.router
}
