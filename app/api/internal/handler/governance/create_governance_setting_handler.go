package governance

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"metaLand/app/api/internal/logic/governance"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
)

// 获取TokenList
func CreateGovernanceSettingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrUpdateGovernanceSettingRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := governance.NewCreateGovernanceSettingLogic(r.Context(), svcCtx)
		resp, err := l.CreateGovernanceSetting(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
