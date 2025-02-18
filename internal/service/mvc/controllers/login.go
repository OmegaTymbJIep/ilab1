package controllers

import (
	"html/template"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type AuthController struct {
	templates *template.Template

	log   *logan.Entry
	model *models.Login
}

func NewAuthController(log *logan.Entry, db data.MainQ, templates *template.Template) *AuthController {
	return &AuthController{
		model:     models.NewLogin(db),
		templates: templates,
		log:       log.WithField("controller", "auth"),
	}
}

func (c *AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := c.templates.ExecuteTemplate(w, views.LoginTemplateName, views.Login{}); err != nil {
		c.log.WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
}
