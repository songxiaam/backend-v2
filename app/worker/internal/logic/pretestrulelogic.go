package logic

import (
	"context"

	"metaLand/app/worker/internal/svc"
	"metaLand/app/worker/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreTestRuleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreTestRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreTestRuleLogic {
	return &PreTestRuleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreTestRuleLogic) PreTestRule(in *pb.PreTestRuleRequest) (*pb.PreTestRuleResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PreTestRuleResponse{}, nil
}
