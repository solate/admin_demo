package menu

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/menu"
	"admin_backend/app/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取所有菜单
func GetAllMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetAllMenuLogic(r.Context(), svcCtx)
		resp, err := l.GetAllMenu()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
