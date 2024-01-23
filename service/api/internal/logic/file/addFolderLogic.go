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

type AddFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFolderLogic {
	return &AddFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFolderLogic) AddFolder(req *types.AddFolderReq) error {
	// 登录用户权限验证（仅限"5-督导老师 42-区域负责人 43-组委会成员 44-组委会主任"可以创建文件夹）
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
	// 添加文件夹（建议role为2或者3，2表示岗前培训及以上可以看，3表示见习队员及以上可以看）
	folder := new(models.Folder)
	folder.FolderName = req.FolderName
	folder.Role = req.Role // 身份代号 -1-待处理人员 0-非在册队员 1-申请队员 2-岗前培训 3-见习队员 4-正式队员 5-督导老师 6-树洞之友 40-普通队员 41-核心队员 42-区域负责人 43-组委会成员 44-组委会主任
	if userInfo.Role < req.Role {
		return response.Error(403, "权限不够")
	}
	_, err = l.svcCtx.FolderModel.Insert(l.ctx, folder)
	if err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success("添加文件夹成功！")
}
