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

type EditFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditFolderLogic {
	return &EditFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditFolderLogic) EditFolder(req *types.EditFolderReq) error {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以修改文件夹信息）
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
	// 修改文件夹信息
	folderInfo, err := l.svcCtx.FolderModel.FindOneNotDel(l.ctx, req.FolderId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return response.Error(100, "当前文件夹不存在")
	default:
		return response.Error(500, err.Error())
	}
	folderInfo.FolderName = req.FolderName
	folderInfo.Role = req.Role
	err = l.svcCtx.FolderModel.Update(l.ctx, folderInfo)
	if err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success("修改文件夹信息成功！")
}
