package file

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"
	"management_system/service/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFolderLogic {
	return &GetFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFolderLogic) GetFolder(req *types.GetFolderReq) (resp *types.GetFolderResp, err error) {
	/*
		获取文件夹列表：
		如果当前用户权限小于文件夹权限，则不能看到该文件夹。
	*/
	// 登录用户权限验证
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(401, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	// 查询文件夹信息
	Folders, err := l.svcCtx.FolderModel.FindAllByUserRole(l.ctx, l.svcCtx.FolderModel.RowBuilder(), req.Page, req.PageSize, userInfo.Role)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	var list []types.Folder
	if len(Folders) > 0 {
		for _, foldermodel := range Folders {
			if foldermodel.Role <= userInfo.Role { // 如果文件夹的权限小于等于用户权限，则用户可以查询到该文件夹
				var typeFolder types.Folder
				_ = copier.Copy(&typeFolder, foldermodel)
				typeFolder.CreateTime = foldermodel.CreateTime.Format("2006-01-02 15:04:05")
				typeFolder.UpdateTime = foldermodel.UpdateTime.Format("2006-01-02 15:04:05")
				list = append(list, typeFolder)
			}
		}
	}
	return &types.GetFolderResp{
		List: list,
	}, nil
}
