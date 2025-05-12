package main

import (
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	Log logx.LogConf
}

func main() {
	var f = flag.String("f", "config.yaml", "config file")
	var c Config

	flag.Parse()
	conf.MustLoad(*f, &c)
	logx.MustSetup(c.Log)

	logx.Info("main thread started!")

	cn := cron.New()

	cn.AddFunc("@every 1s", func() {
		fmt.Println("tick every 1 second")
	})

	cn.Start()

	select {}
}
