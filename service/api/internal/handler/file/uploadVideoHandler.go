package file

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/file"
	"management_system/service/api/internal/svc"
)

func UploadVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := file.NewUploadVideoLogic(r.Context(), svcCtx, r)
		resp, err := l.UploadVideo()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
