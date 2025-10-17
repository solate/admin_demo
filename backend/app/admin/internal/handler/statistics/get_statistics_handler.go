package statistics

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/statistics"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StatisticsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := statistics.NewGetStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.GetStatistics(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
