// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"net/http"

	"admin_backend/internal/logic/task"
	"admin_backend/internal/svc"
	"admin_backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取计划列表
func ListPlanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PlanListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewListPlanLogic(r.Context(), svcCtx)
		resp, err := l.ListPlan(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
