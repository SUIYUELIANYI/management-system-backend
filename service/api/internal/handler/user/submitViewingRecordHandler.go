package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"management_system/service/api/internal/logic/user"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
)

func SubmitViewingRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitViewingRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewSubmitViewingRecordLogic(r.Context(), svcCtx)
		err := l.SubmitViewingRecord(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
