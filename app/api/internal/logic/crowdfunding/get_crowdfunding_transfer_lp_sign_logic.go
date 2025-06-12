package crowdfunding

import (
	"context"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCrowdfundingTransferLpSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetCrowdfundingTeansferLpSign
func NewGetCrowdfundingTransferLpSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCrowdfundingTransferLpSignLogic {
	return &GetCrowdfundingTransferLpSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCrowdfundingTransferLpSignLogic) GetCrowdfundingTransferLpSign(req *types.GetCrowdfundingTransferLpSignRequest) (resp *types.GetCrowdfundingTransferLpSignResponse, err error) {
	// todo: add your logic here and delete this line
	resp = &types.GetCrowdfundingTransferLpSignResponse{
		Sign: "sign",
	}
	return
}
