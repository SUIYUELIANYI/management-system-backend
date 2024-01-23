package rescue

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/rescue"
	"management_system/service/api/internal/svc"
)

func DestroyRescueTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rescue.NewDestroyRescueTokenLogic(r.Context(), svcCtx)
		err := l.DestroyRescueToken()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, nil, err)
	}
}
