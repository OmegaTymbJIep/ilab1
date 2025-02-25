package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

const JWTCookieName = "jwt"

type Auth struct {
	model *models.Auth
}

func NewAuth(model *models.Auth) *Auth {
	return &Auth{
		model: model,
	}
}

func (c *Auth) AuthPage(w http.ResponseWriter, r *http.Request) {
	if customerID, _ := c.verifyJWT(r); customerID != uuid.Nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := Templates(r).ExecuteTemplate(w, views.AuthTemplateName, nil); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewLogin(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	jwt, err := c.model.Login(req)
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			Log(r).WithError(err).Debug("not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}

		if errors.Is(err, models.ErrorInvalidPassword) {
			Unauthorized(w, r, fmt.Errorf("invalid password"))
			return
		}

		InternalError(w, r, fmt.Errorf("failed to login user: %w", err))
		return
	}

	http.SetCookie(w, newJWTCookie(jwt))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Auth) Register(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRegister(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, requests.BadRequest(err)...)
		return
	}

	jwt, err := c.model.Register(req)
	if err != nil {
		if errors.Is(err, models.ErrorEmailOrUsernameTaken) {
			Log(r).WithField("reason", err).Debug("conflict")
			ape.RenderErr(w, problems.Conflict())
			return
		}

		InternalError(w, r, err)
		return
	}

	http.SetCookie(w, newJWTCookie(jwt))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	Log(r).WithField("reason", err).Debug("unauthorized")
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
	return
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, newJWTCookie(&models.JWTWithEat{
		Token:      "",
		Expiration: time.Now().Add(-time.Hour),
	}))
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func newJWTCookie(jwt *models.JWTWithEat) *http.Cookie {
	return &http.Cookie{
		Name:     JWTCookieName,
		Value:    jwt.Token,
		Secure:   false,
		HttpOnly: true,
		Expires:  jwt.Expiration,
		Path:     "/",
	}
}
