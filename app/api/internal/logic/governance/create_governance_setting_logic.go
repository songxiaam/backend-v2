package governance

import (
	"context"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGovernanceSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取TokenList
func NewCreateGovernanceSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGovernanceSettingLogic {
	return &CreateGovernanceSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGovernanceSettingLogic) CreateGovernanceSetting(req *types.CreateOrUpdateGovernanceSettingRequest) (resp *types.CreateGovernanceSettingResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
