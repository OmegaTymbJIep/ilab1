package controllers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	templatesCtxKey
	customerIDCtxKey
)

func CtxLog(entry *logan.Entry) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxTemplates(templates *template.Template) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, templatesCtxKey, templates)
	}
}

func Templates(r *http.Request) *template.Template {
	return r.Context().Value(templatesCtxKey).(*template.Template)
}

func CustomerID(r *http.Request) uuid.UUID {
	return r.Context().Value(customerIDCtxKey).(uuid.UUID)
}
