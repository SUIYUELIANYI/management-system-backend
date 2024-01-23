package rescue

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/rescue"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetRescueDurationThresholdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetRescueDurationThresholdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rescue.NewSetRescueDurationThresholdLogic(r.Context(), svcCtx)
		err := l.SetRescueDurationThreshold(&req)
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, nil, err)
	}
}
