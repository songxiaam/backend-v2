package startup

import (
	"context"

	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStartupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询项目列表
func NewListStartupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStartupsLogic {
	return &ListStartupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListStartupsLogic) ListStartups(req *types.ListStartupsRequest) (resp *types.ListStartupsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
