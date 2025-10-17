package role

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/role"
	"admin_backend/app/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取所有角色列表
func GetAllRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetAllRolesLogic(r.Context(), svcCtx)
		resp, err := l.GetAllRoles()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
