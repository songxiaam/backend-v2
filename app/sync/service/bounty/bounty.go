package bounty

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"                  //以太坊基础库
	ethCommon "github.com/ethereum/go-ethereum/common" //以太坊常用工具
	"github.com/ethereum/go-ethereum/ethclient"        //以太坊客户端
	"github.com/redis/go-redis/v9"                     //redis客户端
	"github.com/zeromicro/go-zero/core/logx"           //日志组件
	"github.com/zeromicro/go-zero/core/threading"      //协程管理
	"math/big"                                         //大整数计算
	"metaLand/app/sync/service/common"
	chainModel "metaLand/data/model/chain" //链数据模型
	"time"
)

// TaskBounty 赏金任务处理器
type TaskBounty struct {
	ctx  *common.ServiceContext   //服务上下文
	info map[uint64]*ContractInfo //链上合约数据，key为chainID
}

func NewTaskBounty(ctx *common.ServiceContext) *TaskBounty {
	return &TaskBounty{
		ctx:  ctx,
		info: make(map[uint64]*ContractInfo),
	}
}

func (t *TaskBounty) HandlerCreateEvent(params []any) {
	name := params[0].(string)
	profile := params[1].(struct {
		Name       string `json:"name"`
		Mode       uint8  `json:"mode"`
		Logo       string `json:"logo"`
		Mission    string `json:"mission"`
		Overview   string `json:"overview"`
		IsValidate bool   `json:"isValidate"`
	})
	comer := params[2].(ethCommon.Address)
	logx.Info(name, comer)
	logx.Info(profile)
	// todo 监听创建bounty create事件，入库

}

func (t *TaskBounty) queryLogs() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for chainID, info := range t.info {
		startBlock := big.NewInt(0)
		bountyABI, err := common.GetABI(info.ABI)
		if err != nil {
			logx.Error(err)
			return
		}
		lastHeigh, err := t.ctx.Redis.Get(ctx, common.GetKey(chainID, info.Address)).Uint64()
		if err != nil && !errors.Is(err, redis.Nil) {
			logx.Error(err)
			continue
		}
		if lastHeigh == 0 {
			receipt, err := info.Client.TransactionReceipt(ctx, ethCommon.HexToHash(info.CreatedHash))
			if err != nil {
				logx.Error(err)
				continue
			}
			startBlock = big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(1))
		} else {
			startBlock = big.NewInt(int64(lastHeigh + 1))
		}
		currentHeight, err := info.Client.BlockNumber(ctx)
		if err != nil {
			logx.Info(err)
			continue
		}
		if big.NewInt(int64(currentHeight)).Cmp(startBlock) <= 0 {
			continue
		}
		endBlock := big.NewInt(0).Add(startBlock, big.NewInt(int64(499)))
		if endBlock.Cmp(big.NewInt(int64(currentHeight))) > 0 {
			endBlock = big.NewInt(int64(currentHeight))
		}

		for {
			logx.Info(fmt.Sprintf("start: %d end: %d current: %d", startBlock.Int64(), endBlock.Int64(), currentHeight))
			logs, err := info.Client.FilterLogs(ctx, ethereum.FilterQuery{
				FromBlock: startBlock,
				ToBlock:   endBlock,
				Addresses: []ethCommon.Address{
					ethCommon.HexToAddress(info.Address),
				},
			})
			if err != nil {
				logx.Info(err)
				break
			}
			for _, l := range logs {
				switch l.Topics[0] {
				case ethCommon.HexToHash(EventCreated):
					params, err := bountyABI.Events["created"].Inputs.UnpackValues(l.Data)
					if err != nil {
						logx.Error(err)
						continue
					}
					t.HandlerCreateEvent(params)
				default:
					logx.Info(l.Topics[0])
				}
			}
			err = t.ctx.Redis.Set(ctx, common.GetKey(chainID, info.Address), endBlock.Uint64(), 0).Err()
			if err != nil {
				logx.Info(err)
				break
			}
			startBlock = big.NewInt(0).Add(endBlock, big.NewInt(int64(1)))
			endBlock = big.NewInt(0).Add(startBlock, big.NewInt(int64(499)))
			if endBlock.Cmp(big.NewInt(int64(currentHeight))) > 0 {
				endBlock = big.NewInt(int64(currentHeight))
			}
			if endBlock.Cmp(startBlock) <= 0 {
				break
			}
		}
	}
}

func (t *TaskBounty) process() {
	chains := make([]chainModel.ChainBasicResponse, 0)
	err := chainModel.GetChainCompleteList(t.ctx.DB, &chains)
	if err != nil {
		logx.Error(err)
	}
	for _, chain := range chains {
		var rpcUrl string
		for _, endpoint := range chain.ChainEndpoints {
			if endpoint.Protocol == 1 {
				rpcUrl = endpoint.URL
			}
		}
		cli, err := ethclient.Dial(rpcUrl)
		if err != nil {
			logx.Error(err)
			continue
		}
		for _, contract := range chain.ChainContracts {
			if contract.Project == 1 {
				_, has := t.info[chain.ChainID]
				if !has {
					t.info[chain.ChainID] = &ContractInfo{
						Address:     contract.Address,
						CreatedHash: contract.CreatedTxHash,
						Client:      cli,
						ABI:         contract.ABI,
					}
				}
			}
		}
	}
	for {
		t.queryLogs()
		time.Sleep(3 * time.Second)
	}
}

func (t *TaskBounty) Start() {
	threading.GoSafe(t.process)
}
