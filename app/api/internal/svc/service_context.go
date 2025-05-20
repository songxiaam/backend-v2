package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"metaLand/app/api/internal/config"
	"metaLand/app/api/internal/middleware"
)

type ServiceContext struct {
	Config             config.Config
	OIDCAuthMiddleware rest.Middleware
	DB                 *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DB.DataSource))
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:             c,
		OIDCAuthMiddleware: middleware.NewOIDCAuthMiddleware().Handle,
		DB:                 db,
	}
}
