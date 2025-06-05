package governance

import (
	"context"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGovernanceSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取token
func NewGetGovernanceSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGovernanceSettingLogic {
	return &GetGovernanceSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGovernanceSettingLogic) GetGovernanceSetting(req *types.GetGovernanceSettingRequest) (resp *types.GetGovernanceSettingResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
