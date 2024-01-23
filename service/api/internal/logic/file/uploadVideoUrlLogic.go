package file

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UploadVideoUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadVideoUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadVideoUrlLogic {
	return &UploadVideoUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadVideoUrlLogic) UploadVideoUrl(req *types.UploadVideoUrlReq) error {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以上传视频）
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(401, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	if userInfo.Role != 5 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	//
	if req.FileName == "" {
		return response.Error(400, "文件名不能为空")
	}
	if req.FolderId == 0 {
		return response.Error(400, "该文件夹不存在")
	}
	if req.Url == "" {
		return response.Error(400, "文件url不能为空")
	}
	//
	file := new(models.File)
	file.FileName = req.FileName
	file.FolderId = req.FolderId
	file.Url = req.Url
	if err := l.svcCtx.FileModel.TransCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.FileModel.Insert(l.ctx, file)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return response.Error(500, err.Error())
	}
	return response.Success("上传文件成功！")
}
