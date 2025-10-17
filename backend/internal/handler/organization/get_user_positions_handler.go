// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"net/http"

	"admin_backend/internal/logic/organization"
	"admin_backend/internal/svc"
	"admin_backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户的岗位列表
func GetUserPositionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserPositionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization.NewGetUserPositionsLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPositions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
