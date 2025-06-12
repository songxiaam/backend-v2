package crowdfunding

import (
	"context"
	"metaLand/data/model/crowdfunding"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCrowdfundingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Crowdfundings列表
func NewGetCrowdfundingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCrowdfundingLogic {
	return &GetCrowdfundingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCrowdfundingLogic) GetCrowdfunding(req *types.GetCrowdfundingRequest) (resp *types.GetCrowdfundingResponse, err error) {

	list, total, err := crowdfunding.GetCrowdfundingList(l.svcCtx.DB, req.Page, req.Size, req.Keyword)
	if err != nil {
		return nil, err
	}

	var resultList = make([]types.CrowdfundingRes, 0)
	for _, item := range list {
		result := types.CrowdfundingRes{
			BaseInfo: types.BaseInfo{
				ID: item.ID,
			},
			Crowdfunding: types.Crowdfunding{
				ChainId:              item.ChainId,
				TxHash:               item.TxHash,
				CrowdfundingContract: item.CrowdfundingContract,
				StartupID:            item.StartupID,
				ComerID:              item.ComerID,
				RaiseGoal:            item.RaiseGoal.String(),
				RaiseBalance:         item.RaiseBalance.String(),
				SellInfo: types.SellInfo{
					SellTokenContract: item.SellTokenContract,
					SellTokenName:     item.SellTokenName,
					SellTokenSymbol:   item.SellTokenSymbol,
					SellTokenDecimals: item.SellTokenDecimals,
					SellTokenSupply:   item.SellTokenSupply.String(),
					SellTokenDeposit:  item.SellTokenDeposit.String(),
					SellTokenBalance:  item.SellTokenBalance.String(),
					MaxSellPercent:    item.MaxSellPercent.String(),
					SellTax:           item.SellTax.String(),
				},
				BuyInfo: types.BuyInfo{
					BuyTokenContract: item.BuyTokenContract,
					BuyTokenName:     item.BuyTokenName,
					BuyTokenSymbol:   item.BuyTokenSymbol,
					BuyTokenDecimals: item.BuyTokenDecimals,
					BuyTokenSupply:   item.BuyTokenSupply.String(),
					BuyPrice:         item.BuyPrice.String(),
					MaxBuyAmount:     item.MaxBuyAmount.String(),
				},
				TeamWallet:  item.TeamWallet,
				SwapPercent: item.SwapPercent.String(),
				StartTime:   item.StartTime.Format(time.DateTime),
				EndTime:     item.EndTime.Format(time.DateTime),
				Poster:      item.Poster,
				Youtube:     item.Youtube,
				Detail:      item.Detail,
				Description: item.Description,
				Status:      int(item.Status),
			},
		}
		resultList = append(resultList, result)
	}

	resp = &types.GetCrowdfundingResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: uint64(total),
		List:  resultList,
	}

	return
}
