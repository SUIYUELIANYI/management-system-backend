package user

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/user"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPhoneNumberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPhoneNumberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetPhoneNumberLogic(r.Context(), svcCtx)
		resp, err := l.GetPhoneNumber(&req)
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
