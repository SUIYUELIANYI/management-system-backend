package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"management-system/service/user/cmd/api/internal/logic/user"
	"management-system/service/user/cmd/api/internal/svc"
	"management-system/service/user/cmd/api/internal/types"
)

func WXMiniAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WXMiniAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewWXMiniAuthLogic(r.Context(), svcCtx)
		resp, err := l.WXMiniAuth(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
