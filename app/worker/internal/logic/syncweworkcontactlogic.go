package logic

import (
	"context"

	"metaLand/app/worker/internal/svc"
	"metaLand/app/worker/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncWeworkContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncWeworkContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncWeworkContactLogic {
	return &SyncWeworkContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncWeworkContactLogic) SyncWeworkContact(in *pb.SyncWeworkContactRequest) (*pb.SyncWeworkContactResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.SyncWeworkContactResponse{}, nil
}
