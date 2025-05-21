package startup

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"math/big"
	"metaLand/app/sync/service/common"
	chainModel "metaLand/data/model/chain"
	"time"
)

type TaskStartup struct {
	ctx  *common.ServiceContext
	info map[uint64]*ContractInfo
}

func NewTaskStartup(ctx *common.ServiceContext) *TaskStartup {
	return &TaskStartup{ctx: ctx, info: make(map[uint64]*ContractInfo)}
}

func (t *TaskStartup) HandleCreateEvent(params []any) {
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
}

func (t *TaskStartup) queryLogs() {
	for chainID, info := range t.info {
		startBlock := big.NewInt(0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		startupABI, err := common.GetABI(info.ABI)
		if err != nil {
			logx.Error(err)
			return
		}

		lastHeigh, err := t.ctx.Redis.Get(ctx, common.GetKey(chainID, info.Address)).Uint64()
		if err != nil && err != redis.Nil {
			logx.Info(err)
			continue
		}

		if lastHeigh == 0 {
			receipt, err := info.Client.TransactionReceipt(ctx, ethCommon.HexToHash(info.CreatedHash))
			if err != nil {
				logx.Info(err)
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
				Addresses: []ethCommon.Address{ethCommon.HexToAddress(info.Address)},
			})

			if err != nil {
				logx.Info(err)
				break
			}

			for _, l := range logs {
				switch l.Topics[0] {
				case ethCommon.HexToHash(EventCreated):
					params, err := startupABI.Events["created"].Inputs.UnpackValues(l.Data)
					if err != nil {
						logx.Error(err)
						continue
					}

					t.HandleCreateEvent(params)
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

func (t *TaskStartup) process() {
	chains := make([]chainModel.ChainBasicResponse, 0)

	err := chainModel.GetChainCompleteList(t.ctx.DB, &chains)
	if err != nil {
		logx.Error(err)
	}

	for _, chain := range chains {
		var rpcurl string

		for _, endpoint := range chain.ChainEndpoints {
			if endpoint.Protocol == 1 {
				rpcurl = endpoint.URL
			}
		}

		cli, err := ethclient.Dial(rpcurl)
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

func (t *TaskStartup) Start() {
	threading.GoSafe(t.process)
}
