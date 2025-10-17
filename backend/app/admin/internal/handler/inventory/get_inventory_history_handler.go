// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package inventory

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/inventory"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取库存操作历史
func GetInventoryHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InventoryListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := inventory.NewGetInventoryHistoryLogic(r.Context(), svcCtx)
		resp, err := l.GetInventoryHistory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
