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

// 创建岗位
func CreatePositionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePositionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization.NewCreatePositionLogic(r.Context(), svcCtx)
		resp, err := l.CreatePosition(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
