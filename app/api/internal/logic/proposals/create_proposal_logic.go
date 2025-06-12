package proposals

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
	"metaLand/data/model/governance"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProposalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create Proposal
func NewCreateProposalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProposalLogic {
	return &CreateProposalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProposalLogic) CreateProposal(req *types.CreateProposalRequest) (resp *types.CreateProposalResponse, err error) {

	// 校验参数
	err = validCreateProposalRequest(req)
	if err != nil {
		return
	}
	// TODO: 从header获取comerId
	//comerId, _ := ctx.Keys[middleware.ComerUinContextKey].(uint64)
	comerId := uint64(1)
	if comerId != req.AuthorComerID {
		return nil, errors.New("invalid comerId")
	}

	req.Status = int(governance.ProposalUpcoming)
	if len(req.Choices) == 0 {
		return nil, errors.New("choices cannot be empty")
	}
	releaseTimestamp, _ := time.Parse(time.DateTime, req.ReleaseTimestamp)
	startTime, _ := time.Parse(time.DateTime, req.StartTime)
	endTime, _ := time.Parse(time.DateTime, req.EndTime)
	proposal := governance.GovernanceProposal{
		StartupID:           req.StartupID,
		AuthorComerID:       req.AuthorComerID,
		AuthorWalletAddress: req.AuthorWalletAddress,
		ChainID:             req.ChainID,
		BlockNumber:         req.BlockNumber,
		ReleaseTimestamp:    releaseTimestamp,
		IPFSHash:            req.IPFSHash,
		Title:               req.Title,
		Description:         req.Description,
		DiscussionLink:      req.DiscussionLink,
		VoteSystem:          req.VoteSystem,
		StartTime:           startTime,
		EndTime:             endTime,
		Status:              req.Status,
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) (er error) {

		if er = governance.CreateProposal(tx, &proposal); er != nil {
			return er
		}
		var choices []*governance.GovernanceChoice
		for _, choice := range req.Choices {
			choices = append(choices, &governance.GovernanceChoice{
				ProposalChoice: governance.ProposalChoice{
					ProposalID: proposal.ID,
					ItemName:   choice.ItemName,
					SeqNum:     choice.SeqNum,
				},
			})
		}
		if er = tx.Create(&choices).Error; er != nil {
			return er
		}
		return nil
	})
	return
}

func validCreateProposalRequest(request *types.CreateProposalRequest) error {
	if request.AuthorComerID == 0 {
		return errors.New("invalid authorComerId")
	}
	if request.StartupID == 0 {
		return errors.New("invalid startupId")
	}
	if strings.TrimSpace(request.AuthorWalletAddress) == "" {
		return errors.New("authorWalletAddress can not be empty")
	}
	if request.ChainID == 0 {
		return errors.New("chainId can not be empty")
	}
	if request.BlockNumber == 0 {
		return errors.New("blockNumber can not be empty")
	}
	if strings.TrimSpace(request.IPFSHash) == "" {
		return errors.New("ipfsHash can not be empty")
	}
	if strings.TrimSpace(request.Title) == "" {
		return errors.New("title can not be empty")
	}
	return nil
}
