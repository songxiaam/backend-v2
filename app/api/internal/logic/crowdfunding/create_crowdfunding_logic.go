package crowdfunding

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"metaLand/data/model/crowdfunding"
	"strings"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCrowdfundingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// CreateCrowdfunding
func NewCreateCrowdfundingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCrowdfundingLogic {
	return &CreateCrowdfundingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCrowdfundingLogic) CreateCrowdfunding(req *types.CreateCrowdfundingRequest) (resp *types.CreateCrowdfundingResponse, err error) {
	err = validRequest(req)
	if err != nil {
		return nil, err
	}

	cModel := crowdfunding.Crowdfunding{
		ChainInfo: crowdfunding.ChainInfo{
			ChainId: req.ChainId,
			TxHash:  req.TxHash,
		},
		SellInfo: crowdfunding.SellInfo{
			SellTokenContract: req.SellTokenContract,
			SellTokenName:     req.SellTokenName,     // 出售代币名称
			SellTokenSymbol:   req.SellTokenSymbol,   // 出售代币符号
			SellTokenDecimals: req.SellTokenDecimals, // 出售代币精度
			SellTokenSupply:   stringToDecimal(req.SellTokenSupply),
			SellTokenDeposit:  stringToDecimal(req.SellTokenDeposit),
			SellTokenBalance:  stringToDecimal(req.SellTokenBalance),
			MaxSellPercent:    stringToDecimal(req.MaxSellPercent),
		},
		BuyInfo: crowdfunding.BuyInfo{
			BuyTokenContract: req.BuyTokenContract,
			BuyTokenName:     req.BuyTokenName,                    // 购买代币名称
			BuyTokenSymbol:   req.BuyTokenSymbol,                  // 购买代币符号
			BuyTokenDecimals: req.BuyTokenDecimals,                // 购买代币精度
			BuyTokenSupply:   stringToDecimal(req.BuyTokenSupply), // 购买代币总量
			BuyPrice:         stringToDecimal(req.BuyPrice),
			MaxBuyAmount:     stringToDecimal(req.MaxBuyAmount),
		},
		CrowdfundingContract: req.CrowdfundingContract,          // 众筹合约地址
		StartupID:            req.StartupID,                     //          `gorm:"column:startup_id" json:"startup_id"`                           // 初创公司ID
		ComerID:              req.ComerID,                       //            `gorm:"column:comer_id" json:"comer_id"`                               // 创始人ID
		RaiseGoal:            stringToDecimal(req.RaiseGoal),    // `gorm:"column:raise_goal;type:decimal(38,18)" json:"raise_goal"`                         // 募资目标
		RaiseBalance:         stringToDecimal(req.RaiseBalance), // 已募资金额
		TeamWallet:           req.TeamWallet,                    // 团队钱包地址
		SwapPercent:          stringToDecimal(req.SwapPercent),  // 兑换百分比
		StartTime:            stringToTime(req.StartTime),       // 开始时间
		EndTime:              stringToTime(req.EndTime),         //             time.Time          `gorm:"column:end_time" json:"end_time"`                // 结束时间
		Poster:               req.Poster,                        //               string             `gorm:"column:poster" json:"poster"`                    // 海报URL
		Youtube:              req.Youtube,                       //             string             `gorm:"column:youtube" json:"youtube"`                  // YouTube链接
		Detail:               req.Detail,                        //              string             `gorm:"column:detail" json:"detail"`                    // 详情URL
		Description:          req.Description,                   //          string             `gorm:"column:description" json:"description"`          // 项目描述
		Status:               crowdfunding.CrowdfundingStatus(req.Status),
	}

	err = crowdfunding.CreateCrowdfunding(l.svcCtx.DB, &cModel)
	if err != nil {
		return nil, err
	}
	return
}

func validRequest(r *types.CreateCrowdfundingRequest) error {
	notEmptyString := []string{r.TeamWallet, r.BuyTokenContract, r.Poster, r.Description, r.TxHash}
	for _, target := range notEmptyString {
		if target == "" || strings.TrimSpace(target) == "" {
			return errors.New("invalid param")
		}
	}

	notEmptyInt := []uint64{r.StartupID, r.StartupID}
	for _, target := range notEmptyInt {
		if target == 0 {
			return errors.New("invalid param")
		}
	}
	decimalValue, err := decimal.NewFromString(r.RaiseGoal)
	if err != nil {
		return err
	}
	if decimalValue.IsZero() {
		return errors.New("invalid param")
	}
	if decimalValue.IsNegative() {
		return errors.New("raise goal must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.BuyPrice)
	if err != nil {
		return err
	}
	if decimalValue.IsZero() {
		return errors.New("invalid param")
	}
	if decimalValue.IsNegative() {
		return errors.New("buy Price must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.MaxBuyAmount)
	if err != nil {
		return err
	}
	if decimalValue.IsZero() {
		return errors.New("invalid param")
	}
	if decimalValue.IsNegative() {
		return errors.New("max Buy Amount must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.MaxSellPercent)
	if err != nil {
		return err
	}
	if decimalValue.IsZero() {
		return errors.New("invalid param")
	}
	if decimalValue.IsNegative() {
		return errors.New("max Sell Percent must be positive number")
	}

	if !strings.HasPrefix(r.BuyTokenContract, "0x") || len(r.BuyTokenContract) > 64 {
		return errors.New("Invalid token contract: " + r.BuyTokenContract)
	}

	if !strings.HasPrefix(r.TeamWallet, "0x") || len(r.TeamWallet) > 64 {
		return errors.New("Invalid team wallet address: " + r.TeamWallet)
	}

	decimalValue, err = decimal.NewFromString(r.SwapPercent)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("swap must be positive number")
	}

	decimalValue, err = decimal.NewFromString(r.SellTax)
	if err != nil {
		return err
	}
	if decimalValue.IsNegative() {
		return errors.New("sell Tax must be positive number")
	}

	startTime, err := time.Parse(time.DateTime, r.StartTime)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(time.DateTime, r.EndTime)
	if err != nil {
		return err
	}

	if !startTime.Before(endTime) {
		return errors.New("start time needs to be before End time")
	}

	status := crowdfunding.CrowdfundingStatus(r.Status)
	isValid := crowdfunding.IsValidCrowdfundingStatus(status)
	if !isValid {
		return errors.New("invalid status")
	}

	return nil
}

func stringToDecimal(str string) decimal.Decimal {
	value, _ := decimal.NewFromString(str)
	return value
}
func stringToTime(str string) time.Time {
	value, _ := time.Parse(time.DateTime, str)
	return value
}
