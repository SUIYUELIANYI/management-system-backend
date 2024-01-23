package application

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/application"
	"management_system/service/api/internal/svc"
)

func GetOwnFormHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewGetOwnFormLogic(r.Context(), svcCtx)
		resp, err := l.GetOwnForm()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
