package exam

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/exam"
	"management_system/service/api/internal/svc"
)

func GetGradesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := exam.NewGetGradesLogic(r.Context(), svcCtx)
		resp, err := l.GetGrades()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		} */
		response.SendResponse(w, r, resp, err)
	}
}
