package exam

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/exam"
	"management_system/service/api/internal/svc"
)

func ImportObjectiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := exam.NewImportObjectiveLogic(r.Context(), svcCtx)
		err := l.ImportObjective()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, nil, err)
	}
}
