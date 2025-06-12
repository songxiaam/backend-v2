package governance

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	//"metaLand/app/api/internal/middleware"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
	"metaLand/data/model/governance"
	"strings"

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
	fmt.Println("CreateGovernanceSetting+++++++++++++")
	// 参数校验
	if req.VoteSymbol == "" || strings.TrimSpace(req.VoteSymbol) == "" {
		return nil, errors.New("vote symbol can not be empty")
	}
	proposalThreshold, err := decimal.NewFromString(req.ProposalThreshold)
	if err != nil {
		return
	}
	proposalValidity, err := decimal.NewFromString(req.ProposalValidity)
	if err != nil {
		return
	}

	if len(req.Strategies) == 0 {
		return nil, errors.New("governance strategies can not be empty")
	}

	// TODO: 检查startup是否存在
	//startup, err := startup.GetStartupById(l.svcCtx.DB, req.StartupId)

	// TODO:获取comerId
	//comerId, _ := l.ctx.Value(middleware.ComerUinContextKey).(uint64)
	//if comerId == 0 {
	//	return nil, errors.New("invalid comerId")
	//}
	comerId := uint64(1)

	// TODO: comerId和startup.comerId比较
	//if startup.ComerID != comerId {
	//	return governanceSetting, errors.New("comer is not founder of startup")
	//}
	// 获取当前已存在的governanceSetting
	mayBeExistedSetting, err := governance.GetGovernanceSetting(l.svcCtx.DB, req.StartupId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) (er error) {
		// 没有查到,新建
		if mayBeExistedSetting.ID == 0 {
			mayBeExistedSetting = governance.GovernanceSetting{
				StartupId:         req.StartupId,
				ComerId:           comerId,
				VoteSymbol:        req.VoteSymbol,
				AllowMember:       req.AllowMember,
				ProposalThreshold: proposalThreshold,
				ProposalValidity:  proposalValidity,
			}
			fmt.Println(mayBeExistedSetting)
			er = governance.CreateGovernanceSetting(tx, &mayBeExistedSetting)
			if er != nil {
				return
			}
			er = createStrategiesOrAdmins(tx, req, &mayBeExistedSetting)
			if er != nil {
				return
			}
		} else {
			// updating
			er = governance.UpdateGovernanceSetting(tx, mayBeExistedSetting.ID, &governance.GovernanceSetting{
				VoteSymbol:        req.VoteSymbol,
				AllowMember:       req.AllowMember,
				ProposalThreshold: proposalThreshold,
				ProposalValidity:  proposalValidity,
			})
			if er != nil {
				return
			}
			// select again
			mayBeExistedSetting, _ = governance.GetGovernanceSetting(tx, req.StartupId)
			er = governance.DeleteAdminsBySettingId(tx, mayBeExistedSetting.ID)
			if er != nil {
				return
			}
			er = governance.DeleteStrategiesBySettingId(tx, mayBeExistedSetting.ID)
			if er != nil {
				return
			}
			er = createStrategiesOrAdmins(tx, req, &mayBeExistedSetting)
			if er != nil {
				return
			}
		}
		return nil
	})

	resp = &types.CreateGovernanceSettingResponse{
		GovernanceSetting: types.GovernanceSetting{
			StartupId:         mayBeExistedSetting.StartupId,
			ComerId:           mayBeExistedSetting.ComerId,
			VoteSymbol:        mayBeExistedSetting.VoteSymbol,
			AllowMember:       mayBeExistedSetting.AllowMember,
			ProposalThreshold: mayBeExistedSetting.ProposalThreshold.String(),
			ProposalValidity:  mayBeExistedSetting.ProposalValidity.String(),
		},
	}
	return resp, err
}

// TODO: 逻辑抽离
func createStrategiesOrAdmins(tx *gorm.DB, request *types.CreateOrUpdateGovernanceSettingRequest, mayBeExistedSetting *governance.GovernanceSetting) (er error) {
	var strategies []*governance.GovernanceStrategy

	if len(request.Strategies) > 0 {
		for _, strategy := range request.Strategies {

			tokenMinBalance, er := decimal.NewFromString(strategy.TokenMinBalance)
			if er != nil {
				return er
			}
			strategies = append(strategies, &governance.GovernanceStrategy{
				SettingId:            mayBeExistedSetting.ID,
				DictValue:            strategy.DictValue,
				StrategyName:         strategy.StrategyName,
				ChainId:              strategy.ChainId,
				TokenContractAddress: strategy.TokenContractAddress,
				VoteSymbol:           strategy.VoteSymbol,
				VoteDecimals:         strategy.VoteDecimals,
				TokenMinBalance:      tokenMinBalance,
			})
		}
	}
	er = governance.CreateGovernanceStrategies(tx, strategies)
	if er != nil {
		return
	}
	var admins []*governance.GovernanceAdmin
	if len(request.Admins) > 0 {
		for _, admin := range request.Admins {
			admins = append(admins, &governance.GovernanceAdmin{
				SettingId:     mayBeExistedSetting.ID,
				WalletAddress: admin.WalletAddress,
			})
		}
		return governance.CreateGovernanceAdmins(tx, admins)
	}
	return
}
