package rescue

import (
	"net/http"

	"management_system/common/response"
	"management_system/service/api/internal/logic/rescue"
	"management_system/service/api/internal/svc"
)

func ReleaseResuceInfoByExcelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rescue.NewReleaseResuceInfoByExcelLogic(r.Context(), svcCtx, r)
		err := l.ReleaseResuceInfoByExcel()
		/* if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		} */
		response.SendResponse(w, r, nil, err)
	}
}
