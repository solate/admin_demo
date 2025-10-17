package dict

import (
	"net/http"

	"admin_backend/app/admin/internal/logic/dict"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建字典数据
func CreateDictItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateDictItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict.NewCreateDictItemLogic(r.Context(), svcCtx)
		resp, err := l.CreateDictItem(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
