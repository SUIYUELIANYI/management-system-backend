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

type GetVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLogic {
	return &GetVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoLogic) GetVideo(req *types.GetVideoReq) (resp *types.GetVideoResp, err error) {
	// 用户权限
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "当前用户不存在")
	default:
		return nil, response.Error(500, err.Error())
	}
	folderInfo, err := l.svcCtx.FolderModel.FindOne(l.ctx, req.FolderId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil,response.Error(100, "当前文件夹不存在")
	default:
		return nil,response.Error(500, err.Error())
	}
	if folderInfo.Role > userInfo.Role {
		return nil,response.Error(403, "权限不够")
	}
	// 查询视频
	Files, err := l.svcCtx.FileModel.FindAllByFolderId(l.ctx, l.svcCtx.FileModel.RowBuilder(), req.Page, req.PageSize, req.FolderId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.File
	if len(Files) > 0 {
		for _, filemodel := range Files {
			var typeFile types.File
			_ = copier.Copy(&typeFile, filemodel)
			typeFile.CreateTime = filemodel.CreateTime.Format("2006-01-02 15:04:05")
			typeFile.UpdateTime = filemodel.UpdateTime.Format("2006-01-02 15:04:05")
			list = append(list, typeFile)
		}
	}
	return &types.GetVideoResp{
		List: list,
	}, nil
}
