package common

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"metaLand/app/sync/service/config"
	"metaLand/app/sync/service/eth"
)

type ServiceContext struct {
	Ctx    context.Context
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Eth    *eth.EthClients
}
