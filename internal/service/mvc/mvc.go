package mvc

import (
	"fmt"
	"html/template"

	"github.com/go-chi/chi/v5"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/omegatymbjiep/ilab1/internal/config"
	"github.com/omegatymbjiep/ilab1/internal/data/postgres"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
)

type MVC struct {
	log *logan.Entry

	auth         *controllers.Auth
	accounts     *controllers.Accounts
	transactions *controllers.Transactions

	templates *template.Template
}

func NewMVC(log *logan.Entry, cfg config.Config) (*MVC, error) {
	db := postgres.NewMainQ(cfg.DB())

	templates, err := views.ReadTemplates(cfg.MVC().TemplatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	authModel, err := models.NewAuth(db, cfg.JWT().SigningKey, cfg.JWT().Expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth model: %w", err)
	}

	return &MVC{
		log:          log,
		auth:         controllers.NewAuth(authModel),
		accounts:     controllers.NewAccounts(models.NewMain(db)),
		transactions: controllers.NewTransactions(models.NewTransactions(db, cfg.ATM().PublicKey)),
		templates:    templates,
	}, nil
}

func (m *MVC) Register(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(
			ape.CtxMiddleware(
				controllers.CtxLog(m.log),
				controllers.CtxTemplates(m.templates),
			),
		)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/", m.auth.AuthPage)
			r.Post("/login", m.auth.Login)
			r.Post("/register", m.auth.Register)
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/home", controllers.Homepage)

			r.With(m.auth.VerifyJWT).Route("/", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.Get("/{account-id}", m.accounts.AccountPage)
				})
				r.Get("/logout", controllers.LogOut)
				r.Get("/", m.accounts.AccountListPage)
			})
		})

		r.With(m.auth.VerifyJWT).Route("/api/v1", func(r chi.Router) {
			r.Route("/transactions", func(r chi.Router) {
				r.Post("/deposit", m.transactions.DepositFunds)
				r.Post("/withdraw", m.transactions.WithdrawFunds)
				r.Post("/transfer", m.transactions.TransferFunds)
			})
			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", m.accounts.CreateAccount)
				r.Delete("/{account-id}", m.accounts.DeleteAccount)
				r.Get("/{account-id}/excel", m.accounts.GenerateAccountExcel)
			})
		})
	})
}
