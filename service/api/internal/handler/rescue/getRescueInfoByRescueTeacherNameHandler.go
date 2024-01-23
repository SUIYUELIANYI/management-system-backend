package rescue

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/rescue"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRescueInfoByRescueTeacherNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRescueInfoByRescueTeacherNameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rescue.NewGetRescueInfoByRescueTeacherNameLogic(r.Context(), svcCtx)
		resp, err := l.GetRescueInfoByRescueTeacherName(&req)
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
