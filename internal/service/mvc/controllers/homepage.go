package controllers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if err := Templates(r).ExecuteTemplate(w, views.HomepageTemplateName, nil); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
