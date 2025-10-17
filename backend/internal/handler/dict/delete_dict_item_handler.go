// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"net/http"

	"admin_backend/internal/logic/dict"
	"admin_backend/internal/svc"
	"admin_backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除字典数据
func DeleteDictItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteDictItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict.NewDeleteDictItemLogic(r.Context(), svcCtx)
		resp, err := l.DeleteDictItem(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
