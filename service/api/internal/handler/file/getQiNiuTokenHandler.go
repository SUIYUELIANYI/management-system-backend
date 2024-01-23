package file

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/file"
	"management_system/service/api/internal/svc"
)

func GetQiNiuTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewGetQiNiuTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetQiNiuToken()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
