package organization

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/organization"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 移除用户岗位
func RemoveUserPositionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveUserPositionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := organization.NewRemoveUserPositionLogic(r.Context(), svcCtx)
		resp, err := l.RemoveUserPosition(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
