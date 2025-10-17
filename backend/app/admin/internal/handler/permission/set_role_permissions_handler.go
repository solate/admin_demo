package permission

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/permission"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 设置角色权限
func SetRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetRolePermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := permission.NewSetRolePermissionsLogic(r.Context(), svcCtx)
		resp, err := l.SetRolePermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
