package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type AuthController struct {
	templates *template.Template
	model     *models.Auth
	log       *logan.Entry
}

func NewAuthController(log *logan.Entry, model *models.Auth, templates *template.Template) *AuthController {
	return &AuthController{
		model:     model,
		templates: templates,
		log:       log.WithField("controller", "auth"),
	}
}

func (c *AuthController) AuthPage(w http.ResponseWriter, _ *http.Request) {
	if err := c.templates.ExecuteTemplate(w, views.AuthTemplateName, views.Auth{}); err != nil {
		c.log.WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewLogin(r)
	if err != nil {
		c.log.WithError(err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	token, err := c.model.Login(req)
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			c.log.WithError(err).Debug("not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}

		if errors.Is(err, models.ErrorInvalidPassword) {
			c.log.WithError(err).Debug("unauthorized")
			ape.RenderErr(w, problems.Unauthorized())
			return
		}

		c.log.WithError(err).Error("failed to login user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO: render response here
	ape.Render(w, token)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRegister(r)
	if err != nil {
		c.log.WithError(err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	jwt, err := c.model.Register(req)
	if err != nil {
		if errors.Is(err, models.ErrorEmailOrUsernameTaken) {
			c.log.WithField("reason", err).Debug("conflict")
			ape.RenderErr(w, problems.Conflict())
			return
		}

		c.log.WithError(err).Error("failed to register new user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO: render response here
	ape.Render(w, jwt)
}
