package proposals

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"metaLand/data/model/governance"
	"strings"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VoteProposalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Vote Proposal
func NewVoteProposalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VoteProposalLogic {
	return &VoteProposalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VoteProposalLogic) VoteProposal(req *types.CreateVoteProposalRequest) (resp *types.CreateVoteProposalResponse, err error) {
	// 参数校验
	if req.IpfsHash == "" || strings.TrimSpace(req.IpfsHash) == "" {
		err = errors.New("invalid ipfsHash")
		return
	}
	if req.ChoiceItemId == 0 {
		err = errors.New("invalid choiceItemId")
		return
	}
	votes, err := decimal.NewFromString(req.Votes)
	if err != nil {
		return
	}
	if votes == decimal.Zero || votes.IsZero() || votes.IsNegative() {
		err = errors.New("illegal votes")
		return
	}

	// TODO: 从header获取comer, 并校验
	comerId := uint64(1)

	proposal, err := governance.GetProposalById(l.svcCtx.DB, req.ProposalId)
	if err != nil {
		return
	}
	if proposal.ID == 0 {
		err = fmt.Errorf("proposal %d does not exist", req.ProposalId)
		return
	}
	if proposal.Status != int(governance.ProposalActive) {
		err = errors.New("invalid proposal status")
		return
	}
	choice, err := governance.GetChoiceByProposalIdAndChoiceId(l.svcCtx.DB, req.ProposalId, req.ChoiceItemId)
	if err != nil {
		return
	}
	if choice.ID == 0 {
		err = errors.New("invalid choice")
		return
	}
	governanceVote := governance.GovernanceVote{
		VoteInfo: governance.VoteInfo{
			ProposalID:         req.ProposalId,
			VoterComerID:       comerId,
			VoterWalletAddress: req.VoterWalletAddress,
			ChoiceItemID:       req.ChoiceItemId,
			ChoiceItemName:     req.ChoiceItemName,
			Votes:              votes,
			IPFSHash:           req.IpfsHash,
		},
	}
	pVoteSystem := governance.VoteSystem(proposal.VoteSystem)
	if pVoteSystem == governance.VoteSystemSingleChoiceVoting || pVoteSystem == governance.VoteSystemBasicVoting {
		vote, err := governance.GetVoteRecordByProposalIdAndComerId(l.svcCtx.DB, req.ProposalId, comerId)
		if err != nil {
			return nil, err
		}
		if vote.ID == 0 {
			err = governance.CreateProposalVote(l.svcCtx.DB, &governanceVote)
			return nil, err
		}
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			if er := governance.DeleteVoteByProposalIdAndVoterComer(tx, req.ProposalId, comerId); er != nil {
				return er
			}
			return governance.CreateProposalVote(tx, &governanceVote)
		})
	} else {
		logx.Errorf("unsupported VoteSystem %v\n", pVoteSystem)
	}
	return
}
