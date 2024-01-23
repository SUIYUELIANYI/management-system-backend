package user

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/user"
	"management_system/service/api/internal/svc"
)

func GetAccessTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetAccessTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetAccessToken()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
