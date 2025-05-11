package logic

import (
	"context"

	"metaLand/app/worker/internal/svc"
	"metaLand/app/worker/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUCESchemasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUCESchemasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUCESchemasLogic {
	return &GetUCESchemasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUCESchemasLogic) GetUCESchemas(in *pb.GetUCESchemasRequest) (*pb.GetUCESchemasResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUCESchemasResponse{}, nil
}
