package crontask

import (
	"context"
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"metaLand/app/sync/service/common"
	"metaLand/app/sync/service/eth"
	"metaLand/data/model/bounty"
	"metaLand/data/model/crowdfunding"
	"metaLand/data/model/transaction"
	"runtime/debug"
	"strings"
	"time"
)

func AddTransactionCornTask(cn *cron.Cron, ctx *common.ServiceContext) {
	fmt.Println("AddTransactionCornTask")
	cn.AddFunc("@every 3s", func() {
		err := GetAllContractAddresses(ctx.DB, ctx.Eth)
		if err != nil {
			logx.Error("GetAllContractAddresses error: %v", err)
		}
	})
}

const RetryThreshold = 3

func Recover() {
	if err := recover(); err != nil {
		s := string(debug.Stack())
		log.Error("recover: err=%v\n stack=%s", err, s)
	}
}

// GetAllContractAddresses  todo need refactor...
func GetAllContractAddresses(db *gorm.DB, eth *eth.EthClients) error {
	fmt.Println("GetAllContractAddresses")
	transactions, err := transaction.GetPendingTransactions(db)
	log.Info("####GET ALL TRANSACTION_BY_STATUS:", transactions)
	if err != nil {
		return err
	}

	for _, tran := range transactions {
		contractAddress, status := getContractAddress(tran.ChainID, tran.TxHash, eth)
		time.Sleep(5 * time.Second)
		log.Info("##### handle transaction ", tran, "the contractAddress is:", contractAddress, "the status is :", status)
		sourceID := tran.SourceID
		switch tran.SourceType {
		case transaction.BountyDepositContractCreated:
			go updateBountyContractAndTransactoinStatus(db, sourceID, status, contractAddress)
		case transaction.CrowdfundingContractCreated:
			go func() {
				defer Recover()

				err := func(onChainStatus int, address string) error {
					if onChainStatus == 1 && strings.TrimSpace(address) != "" {

						entity, err := crowdfunding.GetCrowdfundingById(db, sourceID)
						if err != nil {
							return err
						}
						var startTime = entity.StartTime
						var st = crowdfunding.Upcoming
						if startTime.After(time.Now().Add(-10*time.Second)) && startTime.Before(time.Now().Add(10*time.Second)) {
							st = crowdfunding.Live
						}
						if err := crowdfunding.UpdateCrowdfundingContractAddressAndStatus(db, sourceID, address, st); err != nil {
							return err
						}

						if err := transaction.UpdateTransactionStatusById(db, tran.TransactionId, status); err != nil {
							return err
						}
						return nil
					}
					return postOnChainFailure(db, tran, onChainStatus)
				}(status, contractAddress)
				if err != nil {
					logx.Errorf("#### handle transaction CrowdfundingContractCreated error: %v\n", err)
				}
			}()
		case transaction.CrowdfundingModified:
			go func() {
				defer Recover()

				err := func(onChainStatus int, address string) error {
					if onChainStatus == 1 {

						logx.Infof("#### Crowdfunding modified successfully, then update crowdfunding with history, %d\n", sourceID)

						history, err := crowdfunding.GetIboRateHistoryById(db, sourceID)
						if err != nil {
							return err
						}
						if err := crowdfunding.UpdateCrowdfunding(db, history.CrowdfundingId, crowdfunding.ModifyRequest{
							TransactionHashRequest: crowdfunding.TransactionHashRequest{},
							SwapPercent:            history.SwapPercent,
							BuyPrice:               history.BuyPrice,
							MaxBuyAmount:           history.MaxBuyAmount,
							MaxSellPercent:         history.MaxSellPercent,
							EndTime:                history.EndTime,
						}); err != nil {
							return err
						}

						if err := transaction.UpdateTransactionStatusById(db, tran.TransactionId, status); err != nil {
							return err
						}
						return nil
					}
					return postOnChainFailure(db, tran, onChainStatus)
				}(status, contractAddress)
				if err != nil {
					logx.Errorf("#### handle transaction CrowdfundingContractCreated error: %v\n", err)
				}
			}()
		case transaction.CrowdfundingRemoved:
			go func() {
				defer Recover()

				err := func(onChainStatus int, address string) error {
					if onChainStatus == 1 {

						if err := crowdfunding.UpdateCrowdfundingStatus(db, sourceID, crowdfunding.Ended); err != nil {
							return err
						}

						if err := transaction.UpdateTransactionStatusById(db, tran.TransactionId, status); err != nil {
							return err
						}
						return nil
					}
					return postOnChainFailure(db, tran, onChainStatus)
				}(status, contractAddress)
				if err != nil {
					logx.Errorf("#### handle transaction CrowdfundingContractCreated error: %v\n", err)
				}
			}()
		case transaction.CrowdfundingCancelled:
			go func() {
				defer Recover()

				err := func(onChainStatus int, address string) error {
					if onChainStatus == 1 {

						if err := crowdfunding.UpdateCrowdfundingStatus(db, sourceID, crowdfunding.Cancelled); err != nil {
							return err
						}

						if err := transaction.UpdateTransactionStatusById(db, tran.TransactionId, status); err != nil {
							return err
						}
						return nil
					}
					return postOnChainFailure(db, tran, onChainStatus)
				}(status, contractAddress)
				if err != nil {
					logx.Errorf("#### handle transaction CrowdfundingContractCreated error: %v\n", err)
				}
			}()
		case transaction.CrowdfundingBought, transaction.CrowdfundingSold:
			go func() {
				defer Recover()

				err := func(onChainStatus int, address string) error {
					swap, err := crowdfunding.GetCrowdfundingSwapById(db, sourceID)
					if err != nil {
						logx.Errorf("##### GetCrowdfundingSwapById: %s, %d, %v\n", tran.TxHash, sourceID, err)
						return err
					}
					if err := handleOnChainStateForInvestment(db, onChainStatus, swap, *tran); err != nil {
						return err
					}
					return postOnChainFailure(db, tran, onChainStatus)
				}(status, contractAddress)
				if err != nil {
					logx.Errorf("#### handle transaction CrowdfundingContractCreated error: %v\n", err)
				}
			}()
		default:
			panic(fmt.Sprintf("unsupported source type：%d", tran.SourceType))
		}
	}
	return nil
}

func postOnChainFailure(tx *gorm.DB, tran *transaction.GetTransaction, status int) error {
	logx.Infof("#### post on chain failure %d, %s[SOURCE_TYPE:%d]...\n", tran.TransactionId, tran.TxHash, tran.SourceType)

	if status == 2 && tran.RetryTimes < RetryThreshold {
		logx.Infof("#### will retry %d, %s[SOURCE_TYPE:%d]...\n", tran.TransactionId, tran.TxHash, tran.SourceType)
		return transaction.UpdateTransactionStatusById(tx, tran.TransactionId, 0)
	}
	if status == 2 && tran.RetryTimes >= RetryThreshold-1 {
		if err := transaction.UpdateTransactionStatusById(tx, tran.TransactionId, status); err != nil {
			logx.Error("#### on-chain-failure: %d, %s[SOURCE_TYPE:%d]...\n", tran.TransactionId, tran.TxHash, tran.SourceType)
			return err
		}
		switch tran.SourceType {
		case transaction.CrowdfundingContractCreated:
			return crowdfunding.UpdateCrowdfundingStatus(tx, tran.SourceID, crowdfunding.OnChainFailure)
		case transaction.CrowdfundingModified:
			logx.Error("#### modify crowdfunding failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return nil
		case transaction.CrowdfundingRemoved:
			logx.Error("#### modify crowdfunding failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return nil
		case transaction.CrowdfundingCancelled:
			logx.Error("#### modify crowdfunding failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return nil
		case transaction.CrowdfundingBought:
			logx.Error("#### modify crowdfunding failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return crowdfunding.UpdateCrowdfundingSwapStatus(tx, tran.SourceID, crowdfunding.SwapFailure)
		case transaction.CrowdfundingSold:
			logx.Error("#### modify crowdfunding failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return crowdfunding.UpdateCrowdfundingSwapStatus(tx, tran.SourceID, crowdfunding.SwapFailure)
		default:
			logx.Error("#### unknown failure..%d.%s..%d", tran.TransactionId, tran.TxHash, tran.SourceType)
			return nil
		}
	}

	return nil
}

func getContractAddress(chainID uint64, txHashString string, eth *eth.EthClients) (contractAddress string, status int) {
	txHash := common2.HexToHash(txHashString)
	client, err := eth.GetClient(chainID)
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

func handleOnChainStateForInvestment(db *gorm.DB, onChainStatus int, swap crowdfunding.CrowdfundingSwap, tran transaction.GetTransaction) error {
	sourceID := swap.ID
	txHash := swap.TxHash
	return db.Transaction(func(tx *gorm.DB) error {
		// insert or update crowdfundingInvestor while onChainsStatus is 1
		if onChainStatus == 1 {
			logx.Infof("##### Process on-chain-successful-investment: %s, %d, investAccess %s\n", txHash, sourceID, swap.Access.String())
			// update crowdfunding swap onChainStatus:
			if err := crowdfunding.UpdateCrowdfundingSwapStatus(tx, sourceID, crowdfunding.SwapSuccess); err != nil {
				logx.Errorf("##### UpdateCrowdfundingSwapStatus: %s, %d, %v\n", txHash, sourceID, err)
				return err
			}
			investor, err := crowdfunding.FirstOrCreateInvestor(tx, swap.CrowdfundingID, swap.ComerID)
			if err != nil {
				logx.Errorf("##### FirstOrCreateInvestor: %s, %d, %v\n", txHash, sourceID, err)
				return err
			}

			(&investor).Invest(swap.Access, swap.BuyTokenAmount, swap.SellTokenAmount)
			logx.Infof("##### Update investor --- %v \n", investor)
			if err := crowdfunding.UpdateCrowdfundingInvestor(tx, investor); err != nil {
				logx.Errorf("##### UpdateCrowdfundingInvestor: %s, %d, %v\n", txHash, sourceID, err)
				return err
			}

			if err := crowdfunding.UpdateCrowdfundingRaiseBalance(tx, swap); err != nil {
				logx.Errorf("##### UpdateCrowdfundingRaiseBalance: %s, %d, %v\n", txHash, sourceID, err)
				return err
			}

			if err := transaction.UpdateTransactionStatusById(tx, tran.TransactionId, 1); err != nil {
				return err
			}
			return nil
		} else if onChainStatus == 2 && tran.RetryTimes >= 3 {
			return crowdfunding.UpdateCrowdfundingSwapStatus(db, sourceID, crowdfunding.SwapFailure)
		} else {
			logx.Infof("##### [SwapId:%d, TxHash: %s]Retry .....", swap.ID, swap.TxHash)
		}
		return nil
	})
}

func updateBountyContractAndTransactoinStatus(tx *gorm.DB, bountyID uint64, status int, contractAddress string) {
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			log.Error("recover: err=%v\n stack=%s", err, s)
		}
	}()

	logx.Infof("#####UpdateBountyContractAndTransactoinStatus: bountyId-> %d, status->%d, contractAddress->%s", bountyID, status, contractAddress)
	err := transaction.UpdateTransactionStatus(tx, bountyID, status)
	if err != nil {
		logx.Error(err)
	}

	err = bounty.UpdateBountyDepositContract(tx, bountyID, contractAddress)
	if err != nil {
		logx.Error(err)
	}
	err = bounty.UpdateBountyDepositStatus(tx, bountyID, status)
	if err != nil {
		logx.Error(err)
	}
}
