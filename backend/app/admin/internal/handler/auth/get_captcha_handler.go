package auth

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/auth"
	"admin_backend/app/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取验证码
func GetCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.GetCaptcha()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
