package rescue

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

type GetSignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSignatureLogic {
	return &GetSignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSignatureLogic) GetSignature(req *types.GetSignatureReq) (resp *types.GetSignatureResp, err error) {
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
	// 判断权限
	if userInfo.Role != 3 && userInfo.Role != 4 && userInfo.Role != 5 && userInfo.Role != 40 && userInfo.Role != 41 && userInfo.Role != 42 && userInfo.Role != 43 && userInfo.Role != 44 {
		return nil, response.Error(403, "权限不够")
	}
	// 查看救援对象
	_, err = l.svcCtx.RescueTargetModel.FindOne(l.ctx, req.RescueTargetId)
	switch err {
	case nil:
	case models.ErrNotFound:
		return nil, response.Error(100, "该救援对象不存在！")
	default:
		return nil, response.Error(500, err.Error())
	}

	var list []types.Signature

	signatures, err := l.svcCtx.SignatureModel.FindAllByRescueTargetId(l.ctx, l.svcCtx.SignatureModel.RowBuilder(), req.RescueTargetId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}
	if len(signatures) > 0 {
		for _, signature := range signatures {
			var t types.Signature
			_ = copier.Copy(&t, signature)
			t.CreateTime = signature.CreateTime.Format("2006-01-02 15:04:05")
			t.UpdateTime = signature.UpdateTime.Format("2006-01-02 15:04:05")
			t.RescueTeacherName, err = l.svcCtx.UserModel.FindUsernameById(l.ctx, signature.RescueTeacherId)
			if err != nil {
				return nil, response.Error(500, err.Error())
			}
			t.RescueTeacherRole, err = l.svcCtx.UserModel.FindRoleById(l.ctx, signature.RescueTeacherId)
			if err != nil {
				return nil, response.Error(500, err.Error())
			}
			list = append(list, t)
		}
	}

	return &types.GetSignatureResp{
		List: list,
	}, nil
}
