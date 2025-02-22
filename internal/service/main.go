package service

import (
	"context"
	"log"

	"github.com/omegatymbjiep/ilab1/internal/config"
	apim "github.com/omegatymbjiep/ilab1/internal/service/api"
	mvcm "github.com/omegatymbjiep/ilab1/internal/service/mvc"
)

type Service struct {
}

func newService(cfg config.Config) *Service {
	return &Service{}
}

func Run(ctx context.Context, cfg config.Config) {
	svc := newService(cfg)

	apiLogger := cfg.Log().WithField("service", "api")
	api, err := apim.New(apiLogger, cfg.Listener(), svc)
	if err != nil {
		log.Fatalf("failed to init API: %v", err)
	}

	mvcLogger := cfg.Log().WithField("service", "mvc")
	mvc, err := mvcm.NewMVC(mvcLogger, cfg)
	if err != nil {
		log.Fatalf("failed to init MVC: %v", err)
	}

	mvc.Register(api.Router())

	api.Run(ctx)
}
