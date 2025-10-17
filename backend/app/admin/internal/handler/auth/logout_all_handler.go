// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/auth"
	"admin_backend/app/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户登出（所有设备）
func LogoutAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutAllLogic(r.Context(), svcCtx)
		resp, err := l.LogoutAll()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
