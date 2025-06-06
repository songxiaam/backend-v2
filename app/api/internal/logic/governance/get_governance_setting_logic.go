package governance

import (
	"context"
	"metaLand/data/model/governance"

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

	setting, err := governance.GetGovernanceSetting(l.svcCtx.DB, req.StartupId)
	if err != nil {
		return nil, err
	}
	resp = &types.GetGovernanceSettingResponse{
		GovernanceSetting: types.GovernanceSetting{
			StartupId:         setting.StartupId,
			ComerId:           setting.ComerId,
			VoteSymbol:        setting.VoteSymbol,
			AllowMember:       setting.AllowMember,
			ProposalThreshold: setting.ProposalThreshold.String(),
			ProposalValidity:  setting.ProposalValidity.String(),
		},
		Strategies: []types.GovernanceStrategy{},
		Admins:     []types.GovernanceAdmin{},
	}
	return
}
