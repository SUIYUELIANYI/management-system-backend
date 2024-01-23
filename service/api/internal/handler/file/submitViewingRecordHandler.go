package file

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/file"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SubmitViewingRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitViewingRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewSubmitViewingRecordLogic(r.Context(), svcCtx)
		err := l.SubmitViewingRecord(&req)
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, nil, err)
	}
}
