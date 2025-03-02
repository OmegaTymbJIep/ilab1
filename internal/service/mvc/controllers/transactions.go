package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/responses"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
)

type Transactions struct {
	model *models.Transactions
}

func NewTransactions(model *models.Transactions) *Transactions {
	return &Transactions{
		model: model,
	}
}

func (c *Transactions) DepositFunds(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewDeposit(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	newBalance, err := c.model.DepositFunds(CustomerID(r), req)
	switch {
	case errors.Is(err, models.ErrorInvalidATMSignature):
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(models.ErrorInvalidATMSignature)...)
		return
	case errors.Is(err, models.ErrorAccountNotFound):
		Log(r).WithField("reason", err).Debug("not found")
		ape.RenderErr(w, problems.NotFound())
		return
	case errors.Is(err, models.ErrorATMSignatureNotUnique):
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(models.ErrorATMSignatureNotUnique)...)
		return
	case err != nil:
		InternalError(w, r, fmt.Errorf("failed to deposit funds: %w", err))
	}

	ape.Render(w, responses.NewTransactionResult(newBalance))
}

func (c *Transactions) WithdrawFunds(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewWithdrawal(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	newBalance, err := c.model.WithdrawFunds(CustomerID(r), req)
	switch {
	case errors.Is(err, models.ErrorAccountNotFound):
		Log(r).WithField("reason", err).Debug("not found")
		ape.RenderErr(w, problems.NotFound())
		return
	case errors.Is(err, models.ErrorInsufficientFunds):
		Log(r).WithField("reason", err).Debug("forbidden")
		ape.RenderErr(w, problems.Forbidden())
		return
	case err != nil:
		InternalError(w, r, fmt.Errorf("failed to withdraw funds: %w", err))
		return
	}

	ape.Render(w, responses.NewTransactionResult(newBalance))
}

func (c *Transactions) TransferFunds(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewTransfer(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	newBalance, err := c.model.TransferFunds(CustomerID(r), req)
	if err != nil {
		if errors.Is(err, models.ErrorAccountNotFound) {
			Log(r).WithField("reason", err).Debug("not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}
		if errors.Is(err, models.ErrorInsufficientFunds) {
			Log(r).WithField("reason", err).Debug("forbidden")
			ape.RenderErr(w, problems.Forbidden())
			return
		}
		if errors.Is(err, models.ErrorRecipientNotFound) {
			Log(r).WithField("reason", err).Debug("not found")
			ape.RenderErr(w, notFound("recipient account not found"))
			return
		}
		InternalError(w, r, fmt.Errorf("failed to transfer funds: %w", err))
		return
	}

	ape.Render(w, responses.NewTransactionResult(newBalance))
}

func notFound(message string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusNotFound),
		Status: fmt.Sprintf("%d", http.StatusNotFound),
		Detail: message,
	}
}
