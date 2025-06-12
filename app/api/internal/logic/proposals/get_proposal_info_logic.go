package proposals

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"metaLand/data/model/governance"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProposalInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Proposal Info
func NewGetProposalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProposalInfoLogic {
	return &GetProposalInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProposalInfoLogic) GetProposalInfo(req *types.GetProposalInfoRequest) (resp *types.GetProposalInfoResponse, err error) {
	// todo: add your logic here and delete this line

	proposal, err := governance.GetProposalById(l.svcCtx.DB, req.ProposalId)
	if err != nil {
		return nil, err
	}

	if proposal.ID == 0 {
		return nil, errors.New("Proposal Not Found")
	}

	choices, err := governance.GetProposalChoices(l.svcCtx.DB, req.ProposalId)
	resultChoices := make([]types.ProposalChoiceRes, len(choices))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
		}
	}
	for _, choice := range choices {
		resultChoices = append(resultChoices, types.ProposalChoiceRes{
			ProposalID: choice.ProposalID,
			ItemName:   choice.ItemName,
			SeqNum:     choice.SeqNum,
		})
	}
	// TODO: 查询comer和startUp
	// ...

	resp = &types.GetProposalInfoResponse{
		BaseInfo: types.BaseInfo{
			ID: proposal.ID,
		},
		GovernanceProposal: types.GovernanceProposal{
			StartupID:           proposal.StartupID,
			AuthorComerID:       proposal.AuthorComerID,
			AuthorWalletAddress: proposal.AuthorWalletAddress,
			ChainID:             proposal.ChainID,
			BlockNumber:         proposal.BlockNumber,
			ReleaseTimestamp:    proposal.ReleaseTimestamp.Format(time.DateTime),
			IPFSHash:            proposal.IPFSHash,
			Title:               proposal.Title,
			Description:         proposal.Description,
			DiscussionLink:      proposal.DiscussionLink,
			VoteSystem:          proposal.VoteSystem,
			StartTime:           proposal.StartTime.Format(time.DateTime),
			EndTime:             proposal.EndTime.Format(time.DateTime),
			Status:              proposal.Status,
		},
		Choices: resultChoices,
		Startup: types.StartupRes{},
		Comer:   types.ComerResponse{},
	}

	return
}
