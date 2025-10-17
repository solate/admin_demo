package permission

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/permission"
	"admin_backend/app/admin/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取资源类型列表
func GetResourceTypesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := permission.NewGetResourceTypesLogic(r.Context(), svcCtx)
		resp, err := l.GetResourceTypes()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
