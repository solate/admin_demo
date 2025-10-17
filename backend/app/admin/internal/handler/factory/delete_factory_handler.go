// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package factory

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/factory"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除工厂
func DeleteFactoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteFactoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := factory.NewDeleteFactoryLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFactory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
