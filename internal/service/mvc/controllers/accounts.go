package controllers

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/responses"
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

func (c *Main) CreateAccount(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateAccount(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	account, err := c.model.CreateAccount(CustomerID(r), req)
	if err != nil {
		InternalError(w, r, fmt.Errorf("failed to create account: %w", err))
		return
	}

	ape.Render(w, responses.NewCreateAccount(account))
}

func (c *Main) MainPage(w http.ResponseWriter, r *http.Request) {
	accounts, err := c.model.GetAccounts(CustomerID(r))
	if err != nil {
		InternalError(w, r, fmt.Errorf("failed to get accounts: %w", err))
		return
	}

	viewData := &views.Main{
		Accounts: accounts,
	}

	if err := Templates(r).ExecuteTemplate(w, views.AccountsTemplateName, viewData); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
