package auth

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/auth"
	"admin_backend/app/admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户登出
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx, r)
		resp, err := l.Logout()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
