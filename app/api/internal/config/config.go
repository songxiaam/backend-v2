package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	RPCClientConf zrpc.RpcClientConf
	DB            struct {
		DataSource string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
