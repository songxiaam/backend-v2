package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"metaLand/app/api/internal/config"
	"metaLand/app/api/internal/middleware"
	"metaLand/app/api/model"
	"metaLand/app/worker/metalandworker"
)

type ServiceContext struct {
	Config             config.Config
	OIDCAuthMiddleware rest.Middleware
	ChainModel         model.ChainModel
	StartupModel       model.StartupModel
	RPCClient          metalandworker.MetaLandWorker
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		OIDCAuthMiddleware: middleware.NewOIDCAuthMiddleware().Handle,
		ChainModel:         model.NewChainModel(sqlx.NewMysql(c.DB.DataSource)),
		StartupModel:       model.NewStartupModel(sqlx.NewMysql(c.DB.DataSource)),
		RPCClient:          metalandworker.NewMetaLandWorker(zrpc.MustNewClient(c.RPCClientConf)),
	}
}
