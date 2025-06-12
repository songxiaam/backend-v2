package crowdfunding

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"metaLand/app/api/internal/logic/crowdfunding"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
)

// GetCrowdfundingTeansferLpSign
func GetCrowdfundingTransferLpSignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCrowdfundingTransferLpSignRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := crowdfunding.NewGetCrowdfundingTransferLpSignLogic(r.Context(), svcCtx)
		resp, err := l.GetCrowdfundingTransferLpSign(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
