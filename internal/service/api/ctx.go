package api

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	serviceCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxService(service IService) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, serviceCtxKey, service)
	}
}

func Service(r *http.Request) IService {
	return r.Context().Value(serviceCtxKey).(IService)
}
