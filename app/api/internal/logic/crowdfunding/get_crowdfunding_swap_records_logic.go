package crowdfunding

import (
	"context"
	"metaLand/data/model/crowdfunding"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCrowdfundingSwapRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetCrowdfundingSwapRecords
func NewGetCrowdfundingSwapRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCrowdfundingSwapRecordsLogic {
	return &GetCrowdfundingSwapRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCrowdfundingSwapRecordsLogic) GetCrowdfundingSwapRecords(req *types.GetCrowdfundingSwapRecordsRequest) (resp *types.GetCrowdfundingSwapRecordsResponse, err error) {
	list, err := crowdfunding.QuerySwapListByCrowdfundingId(l.svcCtx.DB, req.CrowdfundingId)
	if err != nil {
		return nil, err
	}

	resultList := make([]types.Swap, len(list))
	for _, item := range list {
		result := types.Swap{
			BaseInfo: types.BaseInfo{
				ID: item.ID,
			},
			ChainId:         item.ChainId,
			TxHash:          item.TxHash,
			Timestamp:       item.Timestamp.Format(time.DateTime),
			Status:          int(item.Status),
			CrowdfundingID:  item.CrowdfundingID,
			ComerID:         item.ComerID,
			Access:          int(item.Access),
			BuyTokenSymbol:  item.BuyTokenSymbol,
			BuyTokenAmount:  item.BuyTokenAmount.String(),
			SellTokenSymbol: item.SellTokenSymbol,
			SellTokenAmount: item.SellTokenAmount.String(),
			Price:           item.Price.String(),
		}
		resultList = append(resultList, result)
	}

	resp = &types.GetCrowdfundingSwapRecordsResponse{
		List: resultList,
	}
	return
}
