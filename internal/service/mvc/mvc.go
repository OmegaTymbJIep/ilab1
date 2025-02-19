package mvc

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/omegatymbjiep/ilab1/internal/config"
	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type MVC struct {
	templatesDir string

	log  *logan.Entry
	auth *controllers.AuthController
}

func NewMVC(log *logan.Entry, db data.MainQ, cfg config.MVCConfig) (*MVC, error) {
	templates, err := views.ReadTemplates(cfg.TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &MVC{
		log:          log,
		templatesDir: cfg.TemplatesDir,
		auth:         controllers.NewAuthController(log, db, templates),
	}, nil
}

func (m *MVC) Register(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", m.auth.AuthPage)
			r.Post("/login", m.auth.Login)
			r.Post("/register", m.auth.Register)
		})
	})
}
