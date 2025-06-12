package crowdfunding

import (
	"context"
	"errors"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"metaLand/data/model/crowdfunding"
	"metaLand/data/model/transaction"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCrowdfundingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// UpdateCrowdfunding
func NewUpdateCrowdfundingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCrowdfundingLogic {
	return &UpdateCrowdfundingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCrowdfundingLogic) UpdateCrowdfunding(req *types.UpdateCrowdfundingRequest) (resp *types.UpdateCrowdfundingResponse, err error) {
	// todo: add your logic here and delete this line

	entity, err := crowdfunding.GetCrowdfundingById(l.svcCtx.DB, req.CrowdfundingId)
	if err != nil {
		return nil, err
	}
	// TODO: 从header中获取comerId
	//comerId := uint64(1)
	//if entity.ComerID != comerId {
	//	return nil, errors.New("current comer is not funder of the crowdfunding")
	//}
	if entity.Status != crowdfunding.Live && entity.Status != crowdfunding.Upcoming {
		return nil, errors.New("the crowdfunding can not be modified")
	}

	err = validUpdateCrowdfundingRequest(req, entity.StartTime)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) (er error) {
		iboRate := crowdfunding.IboRateHistory{
			CrowdfundingId: req.CrowdfundingId,
			EndTime:        stringToTime(req.EndTime),
			MaxBuyAmount:   stringToDecimal(req.MaxBuyAmount),
			MaxSellPercent: stringToDecimal(req.MaxSellPercent),
			BuyPrice:       stringToDecimal(req.BuyPrice),
			SwapPercent:    stringToDecimal(req.SwapPercent),
		}
		er = crowdfunding.CreateIboRateHistory(tx, &iboRate)
		if er != nil {
			return
		}

		sourceType := transaction.CrowdfundingModified

		address, status := l.GetContractAddress(entity.ChainId, entity.TxHash)

		return createTransactionAndQueryContract(tx, &address, status, entity.ChainId, iboRate.ID, req.TxHash, sourceType, func(address string, onChainStatus int) error {
			if onChainStatus == 1 {
				logx.Infof("#### Crowdfunding modified successfully, then update crowdfunding entity, %d\n", req.CrowdfundingId)

				if err := crowdfunding.UpdateCrowdfunding(tx, iboRate.CrowdfundingId, iboRate.MaxSellPercent, iboRate.BuyPrice, iboRate.MaxBuyAmount, iboRate.SwapPercent, iboRate.EndTime); err != nil {
					return err
				}
			}
			return nil
		})
	})

	return
}

func createTransactionAndQueryContract(tx *gorm.DB, address *string, status int, chainId, crowdfundingId uint64, txHash string, sourceType int, callback func(address string, onChainStatus int) error) error {
	if err := transaction.CreateTransaction(tx, &transaction.Transaction{
		ChainID:    chainId,
		TxHash:     txHash,
		TimeStamp:  time.Now(),
		Status:     int(crowdfunding.Pending),
		SourceType: sourceType,
		RetryTimes: 0,
		SourceID:   int64(crowdfundingId),
	}); err != nil {
		return err
	}
	// query contract
	//address, status := getContractAddress(chainId, txHash)
	logx.Infof("#### TX_HASH:%s --> CONTRACT_ADDRESS_AND_ON_CHAIN_STATUS_OF_CREATED_CROWDFUNDING:: %s, %d\n", txHash, address, status)
	if status != 2 {
		if err := transaction.UpdateTransactionStatusWithRetry(tx, crowdfundingId, sourceType, status); err != nil {
			return err
		}
	}
	return callback(*address, status)
}

func validUpdateCrowdfundingRequest(r *types.UpdateCrowdfundingRequest, startTime time.Time) error {
	decimalValue, err := decimal.NewFromString(r.MaxBuyAmount)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("raise goal must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.MaxSellPercent)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("buy Price must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.BuyPrice)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("max Buy Amount must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.SwapPercent)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("max Sell Percent must be positive number")
	}

	endTime, err := time.Parse(time.DateTime, r.EndTime)
	if err != nil {
		return err
	}

	if !startTime.Before(endTime) {
		return errors.New("start time needs to be before End time")
	}

	return nil
}

func (l *UpdateCrowdfundingLogic) GetContractAddress(chainID uint64, txHashString string) (contractAddress string, status int) {
	txHash := common2.HexToHash(txHashString)
	client, err := l.svcCtx.Eth.GetClient(chainID)
	if err != nil {
		logx.Error(err)
		return "", transaction.Failure
	}
	tx, isPending, err := client.RPCClient.TransactionByHash(context.Background(), txHash)
	logx.Infof("#####[TRANSACTION] TransactionByHash:%s ->  isPending-> %v, err-> %v \n", txHash, isPending, err)
	if err != nil {
		logx.Error(err)
		return "", transaction.Failure
	}
	if !isPending {
		receipt, err := client.RPCClient.TransactionReceipt(context.Background(), tx.Hash())
		logx.Infof("#####[TRANSACTION] TransactionReceipt: %s -> Mode: %v,  err-> %v \n", txHash, receipt.Type, err)
		if err != nil {
			logx.Error(err)
			return "", transaction.Failure
		}
		var to *common2.Address
		for _, v := range receipt.Logs {
			if v.Address.String() != "" {
				to = &(v.Address)
				break
			}
		}
		if receipt.Status == transaction.ReceiptFailure {
			return "", transaction.ConfirmFailure
		}

		if to != nil {
			contractAddress = to.String()
		} else {
			logx.Error("#####[TRANSACTION] TransactionReceipt Has Empty Field To: %s \n", txHash)
		}

		return contractAddress, transaction.Success
	}
	return "", transaction.Pending
}
