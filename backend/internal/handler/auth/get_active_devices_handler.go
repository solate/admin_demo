// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"net/http"

	"admin_backend/internal/logic/auth"
	"admin_backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前用户活跃设备数量
func GetActiveDevicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetActiveDevicesLogic(r.Context(), svcCtx)
		resp, err := l.GetActiveDevices()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
