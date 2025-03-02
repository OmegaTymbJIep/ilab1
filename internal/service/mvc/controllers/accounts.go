package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/responses"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type Accounts struct {
	model *models.Accounts
}

func NewAccounts(model *models.Accounts) *Accounts {
	return &Accounts{
		model: model,
	}
}

func (c *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateAccount(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("bad request")
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

func (c *Accounts) AccountListPage(w http.ResponseWriter, r *http.Request) {
	accounts, err := c.model.GetAccountList(CustomerID(r))
	if err != nil {
		InternalError(w, r, fmt.Errorf("failed to get accounts: %w", err))
		return
	}

	viewData := &views.AccountsList{
		Accounts: accounts,
	}

	if err = Templates(r).ExecuteTemplate(w, views.AccountsTemplateName, viewData); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

func (c *Accounts) AccountPage(w http.ResponseWriter, r *http.Request) {
	accountIDRaw := r.PathValue("account-id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(fmt.Errorf("invalid account id"))...)
		return
	}

	account, err := c.model.GetAccount(CustomerID(r), accountID)
	if err != nil {
		if errors.Is(err, models.ErrorAccountNotFound) {
			ape.RenderErr(w, problems.NotFound())
			return
		}

		InternalError(w, r, fmt.Errorf("failed to get account: %w", err))
		return
	}

	transactions, err := c.model.GetAccountTransactions(CustomerID(r), accountID)
	if err != nil {
		InternalError(w, r, fmt.Errorf("failed to get transactions: %w", err))
		return
	}

	viewData := &views.Account{
		Account:      account,
		Transactions: transactions,
	}

	if err := Templates(r).ExecuteTemplate(w, views.AccountTemplateName, viewData); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
