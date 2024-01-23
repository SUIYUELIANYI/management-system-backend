package file

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditVideoLogic {
	return &EditVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditVideoLogic) EditVideo(req *types.EditVideoReq) error {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以修改视频信息）
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前用户不存在")
	default:
		return response.Error(500, err.Error())
	}
	if userInfo.Role != 5 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return response.Error(403, "权限不够")
	}
	// 修改视频文件信息
	vedio, err := l.svcCtx.FileModel.FindOneNotDel(l.ctx, req.VideoId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "该视频不存在或已删除")
	default:
		return response.Error(500, err.Error())
	}
	vedio.FolderId = req.FolderId
	vedio.FileName = req.FileName
	err = l.svcCtx.FileModel.Update(l.ctx, vedio)
	if err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success("修改视频信息成功！")
}
