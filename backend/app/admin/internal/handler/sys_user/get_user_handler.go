package sys_user

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/sys_user"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户详情
func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sys_user.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}