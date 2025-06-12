package crowdfunding

import (
	"context"
	"metaLand/data/model/crowdfunding"
	"metaLand/data/model/startup"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCrowdfundingInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取CrowdfundingInfo
func NewGetCrowdfundingInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCrowdfundingInfoLogic {
	return &GetCrowdfundingInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCrowdfundingInfoLogic) GetCrowdfundingInfo(req *types.GetCrowdfundingInfoRequest) (resp *types.GetCrowdfundingInfoResponse, err error) {

	entity, err := crowdfunding.GetCrowdfundingById(l.svcCtx.DB, req.CrowdfundingId)
	if err != nil {
		return nil, err
	}

	investor, err := crowdfunding.SelectInvestorByCrowdfundingIdAndComerId(l.svcCtx.DB, entity.ID, entity.ComerID)
	if err != nil {
		return nil, err
	}

	investorCount, err := crowdfunding.CountByCrowdfundingId(l.svcCtx.DB, req.CrowdfundingId)
	if err != nil {
		return nil, err
	}

	startupInfo, err := startup.GetStartupInfo(l.svcCtx.DB, &entity.StartupID)
	if err != nil {
		return nil, err
	}

	swaps, err := crowdfunding.QuerySwapListByCrowdfundingId(l.svcCtx.DB, req.CrowdfundingId)
	if err != nil {
		return nil, err
	}
	resultSwaps := make([]types.Swap, len(swaps))
	for _, swap := range swaps {
		result := types.Swap{
			BaseInfo: types.BaseInfo{
				ID: swap.ID,
			},
			ChainId:         swap.ChainId,
			TxHash:          swap.TxHash,
			Timestamp:       swap.Timestamp.String(),
			Status:          int(swap.Status),
			CrowdfundingID:  swap.CrowdfundingID,
			ComerID:         swap.ComerID,
			Access:          int(swap.Access),
			BuyTokenSymbol:  swap.BuyTokenSymbol,
			BuyTokenAmount:  swap.BuyTokenAmount.String(),
			SellTokenSymbol: swap.SellTokenSymbol,
			SellTokenAmount: swap.SellTokenAmount.String(),
			Price:           swap.Price.String(),
		}
		resultSwaps = append(resultSwaps, result)
	}

	resp = &types.GetCrowdfundingInfoResponse{
		CrowdfundingRes: types.CrowdfundingRes{
			BaseInfo: types.BaseInfo{
				ID: entity.ID,
			},
			Crowdfunding: types.Crowdfunding{},
		},
		Investor: types.Investor{
			BaseInfo: types.BaseInfo{
				ID: investor.ID,
			},
			ButTokenBalance:  investor.BuyTokenBalance.String(),
			BuyTokenTotal:    investor.BuyTokenTotal.String(),
			ComerId:          investor.ComerId,
			CrowdfundingId:   investor.CrowdfundingId,
			SellTokenBalance: investor.SellTokenBalance.String(),
			SellTokenTotal:   investor.SellTokenTotal.String(),
		},
		Investors: uint64(investorCount),
		Startup: types.Startup{
			ComerID:              startupInfo.ComerID,
			Name:                 startupInfo.Name,
			Mode:                 startupInfo.Mode,
			Logo:                 startupInfo.Logo,
			Cover:                startupInfo.Cover,
			Mission:              startupInfo.Mission,
			TokenContractAddress: startupInfo.TokenContractAddress,
			Overview:             startupInfo.Overview,
			ChainID:              entity.ChainId,
			TxHash:               startupInfo.TxHash,
			OnChain:              startupInfo.OnChain,
			KYC:                  startupInfo.KYC,
			ContractAudit:        startupInfo.ContractAudit,
		},
		Swaps: resultSwaps,
	}
	return
}
