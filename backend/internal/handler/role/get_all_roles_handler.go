// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"admin_backend/internal/logic/role"
	"admin_backend/internal/svc"
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
