package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type Main struct {
	model *models.Main
}

func NewMain(model *models.Main) *Main {
	return &Main{
		model: model,
	}
}

func (c *Main) MainPage(w http.ResponseWriter, r *http.Request) {
	customerID := CustomerID(r)
	if customerID == uuid.Nil {
		Unauthorized(w, r, fmt.Errorf("customer id not found"))
		return
	}

	viewData := new(views.Main)

	if err := Templates(r).ExecuteTemplate(w, views.MainTemplateName, viewData); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
