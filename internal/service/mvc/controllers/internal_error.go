package controllers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Log(r).WithError(err).Error("internal error")

	if err := Templates(r).ExecuteTemplate(w, views.InternalErrorTemplateName, nil); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	return
}
