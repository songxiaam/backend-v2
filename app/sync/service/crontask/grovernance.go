package crontask

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"metaLand/app/sync/service/common"
	"metaLand/data/model/governance"
)

func AddGrovernanceCornTask(cn *cron.Cron, ctx *common.ServiceContext) {
	fmt.Println("AddGrovernanceCornTask")
	cn.AddFunc("@every 5s", func() {
		err := ActiveProposalStatusSchedule(ctx.DB)
		if err != nil {
			logx.Errorf("#### ActiveProposalStatusSchedule: %v\n", err)
		}
	})
	cn.AddFunc("@every 1s", func() {
		err := EndProposalStatusSchedule(ctx.DB)
		if err != nil {
			logx.Errorf("#### EndProposalStatusSchedule: %v\\n", err)
		}
	})
}

func ActiveProposalStatusSchedule(db *gorm.DB) error {
	fmt.Println("ActiveProposalStatusSchedule")
	list, err := governance.SelectToBeStartedProposalListWithin1Min(db)
	if err != nil {
		//log.Errorf("#### update proposal active status: %v\n", err)
		return err
	}
	if len(list) > 0 {
		for _, c := range list {
			logx.Infof("#### update proposal active status: %d\n", c.ID)
			err = governance.UpdateProposalStatus(db, c.ID, governance.ProposalActive)
			if err != nil {
				//log.Infof("#### update proposal active status: %v\n", err)
				return err
			}
		}
	}
	return nil

}

func EndProposalStatusSchedule(db *gorm.DB) error {
	fmt.Println("EndProposalStatusSchedule")
	list, err := governance.SelectToEndedProposalListWithin1Min(db)
	if err != nil {
		//log.Errorf("#### update proposal ended status: %v\n", err)
		return err
	}

	if len(list) > 0 {
		for _, c := range list {

			// TODO:获取setting
			//_, err := getStartupGovernanceSetting(c.StartupID, db)
			//if err != nil {
			//	return err
			//}

			// TODO: 获取TotalVotes
			//detail, err := getProposal(c.ID, db)
			//if err != nil {
			//	return err
			//}
			status := governance.ProposalEnded
			//if detail.TotalVotes.LessThan(setting.GovernanceSetting.ProposalValidity) {
			//	status = governance.ProposalInvalid
			//}

			err = governance.UpdateProposalStatus(db, c.ID, status)
			if err != nil {
				//log.Infof("#### update proposal %v  status: %v\n", status, err)
				return err
			}
		}
	}
	return nil
}

// TODO: logic中有相似逻辑,后续抽离到公共方法

//func getStartupGovernanceSetting(startupId uint64, db *gorm.DB) (detail *types.CreateGovernanceSettingResponse, err error) {
//	setting, err := governance.GetGovernanceSetting(db, startupId)
//	if err != nil {
//		return detail, err
//	}
//
//	//TODO: 获取startup
//	//st, err := startupInfo.GetStartupById(db, startupId)
//	//if err != nil {
//	//	return detail, err
//	//}
//	startupInfo := startup.Startup{
//		ComerID: uint64(1),
//	}
//
//	comerInfo, err := comer.FindComer(db, startupInfo.ComerID)
//	if err != nil {
//		return nil, err
//	}
//	// if setting.ID == 0 , create default governance-setting for startupInfo
//	if setting.ID == 0 {
//		setting, err = createStartupGovernanceSetting(startupInfo.ComerID, startupId, types.CreateOrUpdateGovernanceSettingRequest{
//			StartupId: startupId,
//			SettingRequest: types.SettingRequest{
//				VoteSymbol:        "vote",
//				AllowMember:       true,
//				ProposalThreshold: "0",
//				ProposalValidity:  "0",
//			},
//			Strategies: []types.StrategyRequest{
//				// ChainId of ethereum
//				// TODO: 配置chainId
//				{DictValue: "ticket", StrategyName: "ticket", VoteSymbol: "", ChainId: 1},
//			},
//			Admins: []types.AdminRequest{
//				// default admin is the funder of startupInfo
//				{WalletAddress: comerInfo.Address},
//			},
//		})
//		if err != nil {
//			return nil, err
//		}
//	}
//	_, err = governance.GetGovernanceStrategies(db, setting.ID)
//	if err != nil {
//		return nil, err
//	}
//	admins, err := governance.GetGovernanceAdmins(db, setting.ID)
//	if err != nil {
//		return nil, err
//	}
//	if len(admins) == 0 {
//		governanceAdmins := []*governance.GovernanceAdmin{{SettingId: setting.ID, WalletAddress: comerInfo.Address}}
//		if err := governance.CreateGovernanceAdmins(db, governanceAdmins); err != nil {
//			return detail, err
//		}
//		admins = append(admins, governanceAdmins...)
//	}
//
//	detail = &types.CreateGovernanceSettingResponse{
//		GovernanceSetting: types.GovernanceSetting{
//			BaseInfo: types.BaseInfo{
//				ID: setting.ID,
//			},
//			StartupId:         setting.ID,
//			ComerId:           setting.ComerId,
//			VoteSymbol:        setting.VoteSymbol,
//			AllowMember:       setting.AllowMember,
//			ProposalThreshold: setting.ProposalThreshold.String(),
//			ProposalValidity:  setting.ProposalValidity.String(),
//		},
//	}
//
//	return detail, nil
//}

//func createStartupGovernanceSetting(comerId, startupId uint64, request types.CreateOrUpdateGovernanceSettingRequest) (setting governance.GovernanceSetting, err error) {
//	startup, err := startupModel.GetStartupById(mysql.DB, startupId)
//	if err != nil {
//		return setting, err
//	}
//	if startup.ComerID != comerId {
//		return setting, errors.New("comer is not founder of startup")
//	}
//	if request.VoteSymbol == "" || strings.TrimSpace(request.VoteSymbol) == "" {
//		return setting, errors.New("vote symbol can not be empty")
//	}
//	if len(request.Strategies) == 0 {
//		return setting, errors.New("governance strategies can not be empty")
//	}
//	mayBeExistedSetting, err := governance.GetGovernanceSetting(mysql.DB, startupId)
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return setting, err
//	}
//
//	err = mysql.DB.Transaction(func(tx *gorm.DB) (er error) {
//		if mayBeExistedSetting.ID == 0 {
//			mayBeExistedSetting = governance.GovernanceSetting{
//				StartupId:         startupId,
//				ComerId:           comerId,
//				VoteSymbol:        request.VoteSymbol,
//				AllowMember:       request.AllowMember,
//				ProposalThreshold: request.ProposalThreshold,
//				ProposalValidity:  request.ProposalValidity,
//			}
//			er = governance.CreateGovernanceSetting(tx, &mayBeExistedSetting)
//			if er != nil {
//				return
//			}
//			return createStrategiesOrAdmins(tx, request, mayBeExistedSetting)
//		}
//		// updating
//		er = governance.UpdateGovernanceSetting(tx, mayBeExistedSetting.ID, &governance.GovernanceSetting{
//			VoteSymbol:        request.VoteSymbol,
//			AllowMember:       request.AllowMember,
//			ProposalThreshold: request.ProposalThreshold,
//			ProposalValidity:  request.ProposalValidity,
//		})
//		if er != nil {
//			return
//		}
//		// select again
//		mayBeExistedSetting, _ = governance.GetGovernanceSetting(tx, startupId)
//		er = governance.DeleteAdminsBySettingId(tx, mayBeExistedSetting.ID)
//		if er != nil {
//			return
//		}
//		er = governance.DeleteStrategiesBySettingId(tx, mayBeExistedSetting.ID)
//		if er != nil {
//			return
//		}
//		if err := createStrategiesOrAdmins(tx, request, mayBeExistedSetting); err != nil {
//			return err
//		}
//		return nil
//	})
//	return mayBeExistedSetting, err
//}
//
//func getProposal(proposalId uint64, db *gorm.DB) (detail governance.ProposalDetail, err error) {
//
//	strategies := setting.Strategies
//	admins := setting.Admins
//	voteResult, err := GetProposalCrtVoteResult(proposalId)
//	if err != nil {
//		return
//	}
//	choices, err := governance.GetProposalChoices(mysql.DB, proposalId)
//	if err != nil {
//		return
//	}
//	dictModel, err := dict.SelectByDictTypeAndLabel(mysql.DB, "voteSystem", publicInfo.VoteSystem)
//	if err != nil {
//		return
//	}
//	detail = governance.ProposalDetail{
//		ProposalPublicInfo:        publicInfo,
//		VoteSystemId:              dictModel.ID,
//		Strategies:                strategies,
//		Admins:                    admins,
//		CurrentProposalVoteResult: voteResult,
//		Choices:                   choices,
//	}
//	return
//}
//
//func getTotalVotes(proposalId uint64, db *gorm.DB) *decimal.Decimal {
//	voteResult, err := GetProposalCrtVoteResult(proposalId)
//}
