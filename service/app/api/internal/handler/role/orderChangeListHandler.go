package role

import (
	"net/http"

	"management-system/response"
	"management-system/service/app/api/internal/logic/role"
	"management-system/service/app/api/internal/svc"
	"management-system/service/app/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OrderChangeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleChangeInfoListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewOrderChangeListLogic(r.Context(), svcCtx)
		resp, err := l.OrderChangeList(&req)
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err) // 错误处理 3.13
	}
}
