package proposals

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
	"metaLand/data/model/governance"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProposalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete Proposal
func NewDeleteProposalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProposalLogic {
	return &DeleteProposalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProposalLogic) DeleteProposal(req *types.DeleteProposalRequest) (resp *types.DeleteProposalResponse, err error) {
	// todo: add your logic here and delete this line

	// TODO: 从Header中获取comer, 校验权限
	comerId := uint64(1)
	//findComer, err := comer.FindComer(l.svcCtx.DB, comerId)
	//if err != nil {
	//	return
	//}
	//if findComer.ID == 0 {
	//	err = errors.New(fmt.Sprintf("invalid comer %d", comerId))
	//	return
	//}
	//
	//if findComer.Address == "" {
	//	err = errors.New(fmt.Sprintf("invalid comer %d without walletAddress", comerId))
	//	return
	//}

	proposal, err := governance.GetProposalById(l.svcCtx.DB, req.ProposalId)
	if err != nil {
		return
	}
	if proposal.ID == 0 {
		err = errors.New(fmt.Sprintf("proposal %d does not exist", req.ProposalId))
		return
	}

	admins, err := governance.GetGovernanceAdminsByStartupId(l.svcCtx.DB, proposal.StartupID)
	if err != nil {
		return
	}
	var can bool
	if len(admins) == 0 {
		can = comerId == proposal.AuthorComerID
	} else {
		//for _, admin := range admins {
		//	can = admin.WalletAddress == findComer.Address
		//	if can {
		//		break
		//	}
		//}
	}
	if !can {
		err = errors.New(fmt.Sprintf("cannot delete proposal"))
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) (er error) {
		er = governance.DeleteProposal(l.svcCtx.DB, proposal.ID)
		if er != nil {
			return
		}
		er = governance.DeleteProposalChoices(l.svcCtx.DB, proposal.ID)
		if er != nil {
			return
		}
		return
	})
	return
}
