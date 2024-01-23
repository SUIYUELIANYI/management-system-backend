package rescue

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/rescue"
	"management_system/service/api/internal/svc"
)

func AreaRescueFrequencyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rescue.NewAreaRescueFrequencyLogic(r.Context(), svcCtx)
		resp, err := l.AreaRescueFrequency()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
