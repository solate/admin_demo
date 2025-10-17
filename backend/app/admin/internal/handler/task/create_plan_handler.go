package task

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/task"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建计划
func CreatePlanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePlanReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewCreatePlanLogic(r.Context(), svcCtx)
		resp, err := l.CreatePlan(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
