package crontask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"metaLand/app/sync/service/common"
	"metaLand/app/sync/service/eth"
	"metaLand/data/model/chain"
)

type Task struct {
	ctx *common.ServiceContext
}

func NewTask(ctx *common.ServiceContext) *Task {
	ethClients := eth.NewEthClients()
	var chains []chain.ChainBasicResponse
	err := chain.GetChainCompleteList(ctx.DB, &chains)
	if err != nil {
		logx.Errorf("GetChainCompleteList error: %v", err)
		panic(err)
	}
	ethClients.Start(&chains)
	ctx.Eth = ethClients
	return &Task{ctx: ctx}
}

func (t *Task) process() {
	cn := cron.New()

	cn.AddFunc("@every 5m", func() {
		fmt.Println("cron working")
	})

	// 添加 crowdfunding 相关定时任务
	AddCrowdfundingCornTask(cn, t.ctx)
	// 添加Grovernance相关定时任务
	AddGrovernanceCornTask(cn, t.ctx)
	// 添加AddTransactionCornTask
	AddTransactionCornTask(cn, t.ctx)

	cn.Start()
}

func (t *Task) Start() {
	threading.GoSafe(t.process)
}
