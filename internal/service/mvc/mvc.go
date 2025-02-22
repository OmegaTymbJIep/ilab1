package mvc

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/omegatymbjiep/ilab1/internal/config"
	"github.com/omegatymbjiep/ilab1/internal/data/postgres"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type MVC struct {
	log  *logan.Entry
	auth *controllers.AuthController
}

func NewMVC(log *logan.Entry, cfg config.Config) (*MVC, error) {
	db := postgres.NewMainQ(cfg.DB())

	templates, err := views.ReadTemplates(cfg.MVC().TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	authModel, err := models.NewAuth(db, cfg.JWT().SigningKey)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth model: %w", err)
	}

	return &MVC{
		log:  log,
		auth: controllers.NewAuthController(log, authModel, templates),
	}, nil
}

func (m *MVC) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", m.auth.AuthPage)
			r.Post("/login", m.auth.Login)
			r.Post("/register", m.auth.Register)
		})

		//r.With(jwtx).Route("/dashboard", func(r chi.Router) {
		//
		//})
	})
}
