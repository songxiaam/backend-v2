package proposals

import (
	"context"
	"metaLand/data/model/governance"
	"time"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProposalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Proposal
func NewGetProposalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProposalsLogic {
	return &GetProposalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProposalsLogic) GetProposals(req *types.GetProposalsRequest) (resp *types.GetProposalsResponse, err error) {
	proposals, total, err := governance.GetProposalList(l.svcCtx.DB, req.Page, req.Size, req.Keyword)
	if err != nil {
		return nil, err
	}
	if len(proposals) == 0 {
		return &types.GetProposalsResponse{
			Page:  req.Page,
			Size:  req.Size,
			Total: 0,
		}, nil
	}

	startupIds := make([]uint64, 0, len(proposals))
	authorComerIDs := make([]uint64, 0, len(proposals))
	for _, proposal := range proposals {
		startupIds = append(startupIds, proposal.StartupID)
		authorComerIDs = append(authorComerIDs, proposal.AuthorComerID)
	}

	// TODO: 填充 StartUp和Comer
	// 根据startUpIds查询StartUpList
	// 根据authorComerIDs查询ComerList

	proposalList := make([]types.GovernanceProposalInfo, len(proposals))
	for _, proposal := range proposals {
		item := types.GovernanceProposalInfo{
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

			//Startup:
			//Comer:
		}
		proposalList = append(proposalList, item)
	}

	resp = &types.GetProposalsResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
		List:  proposalList,
	}

	return
}
