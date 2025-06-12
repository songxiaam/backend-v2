package governance

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"metaLand/app/api/internal/logic/governance"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
)

// 获取token
func GetGovernanceSettingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGovernanceSettingRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := governance.NewGetGovernanceSettingLogic(r.Context(), svcCtx)
		resp, err := l.GetGovernanceSetting(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
