package sys_user

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/sys_user"
	"admin_backend/app/admin/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前用户信息
func GetCurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sys_user.NewGetCurrentUserLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
