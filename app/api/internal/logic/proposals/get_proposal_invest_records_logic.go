package proposals

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
	"metaLand/data/model/governance"
)

type GetProposalInvestRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Proposal Invest Records
func NewGetProposalInvestRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProposalInvestRecordsLogic {
	return &GetProposalInvestRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProposalInvestRecordsLogic) GetProposalInvestRecords(req *types.GetProposalInvestRecordsRequest) (resp *types.GetProposalInvestRecordsResponse, err error) {
	// todo: add your logic here and delete this line

	var total int64
	list, total, err := governance.GetVoteRecordsByProposalId(l.svcCtx.DB, req.ProposalId, req.Page, req.Size)

	resultList := make([]types.GovernanceVote, len(list))
	comerIds := make([]uint64, 0)
	for _, vote := range list {

		item := types.GovernanceVote{
			BaseInfo: types.BaseInfo{
				ID: vote.ID,
			},
			ProposalId:         vote.ProposalID,
			VoterComerId:       vote.VoterComerID,
			VoterWalletAddress: vote.VoterWalletAddress,
			ChoiceItemId:       vote.ChoiceItemID,
			ChoiceItemName:     vote.ChoiceItemName,
			Votes:              vote.Votes.String(),
			IPFSHash:           vote.IPFSHash,
			//Comer              ComerResponse `json:"comer"`
		}
		comerIds = append(comerIds, vote.VoterComerID)
		resultList = append(resultList, item)
	}
	// TODO: 根据comerIds获取comers, 赋值给resultList
	//for i := range resultList {
	//	resultList[i].Comer =
	//}

	resp = &types.GetProposalInvestRecordsResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
		List:  resultList,
	}
	return
}
