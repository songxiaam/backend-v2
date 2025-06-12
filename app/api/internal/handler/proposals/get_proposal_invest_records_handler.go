package proposals

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"metaLand/app/api/internal/logic/proposals"
	"metaLand/app/api/internal/svc"
	"metaLand/app/api/internal/types"
)

// Get Proposal Invest Records
func GetProposalInvestRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetProposalInvestRecordsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := proposals.NewGetProposalInvestRecordsLogic(r.Context(), svcCtx)
		resp, err := l.GetProposalInvestRecords(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
