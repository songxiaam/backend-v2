package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"metaLand/app/api/internal/config"
	"metaLand/app/api/internal/middleware"
)

type ServiceContext struct {
	Config             config.Config
	OIDCAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		OIDCAuthMiddleware: middleware.NewOIDCAuthMiddleware().Handle,
	}
}
