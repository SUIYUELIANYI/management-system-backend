package file

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"management_system/common/qiniu"
	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UploadVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadVideoLogic {
	return &UploadVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadVideoLogic) UploadVideo() (resp *types.UploadVideoResp, err error) {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以上传视频）
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(401, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	if userInfo.Role != 5 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 读取文件和相关信息
	_, fileHeader, err := l.r.FormFile("file")
	fileName := l.r.FormValue("fileName")
	folder := l.r.FormValue("folderId")
	if err != nil {
		return nil, response.Error(400, "文件上传错误")
	}
	if fileName == "" || folder == "" {
		return nil, response.Error(400, "文件名或文件夹不能为空")
	}

	file := new(models.File)
	var url string
	var videoId int64
	if err := l.svcCtx.FileModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		qiniu.Load()
		url, err = qiniu.UploadQiniu(fileHeader)
		if err != nil {
			return err
		}
		file.Url = url
		file.FileName = fileName
		file.FolderId, err = strconv.ParseInt(folder, 10, 64)
		if err != nil {
			return err
		}
		result, err := l.svcCtx.FileModel.Insert(l.ctx, file)
		if err != nil {
			return err
		}
		videoId, _ = result.LastInsertId()
		return nil
	}); err != nil {
		return nil, response.Error(500, err.Error())
	}

	return &types.UploadVideoResp{
		VideoId: videoId,
		Url:     url,
	}, nil
}
