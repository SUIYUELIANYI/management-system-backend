package user

import (
	"context"
	"time"

	"management-system/common/utils"
	"management-system/common/xerr"
	"management-system/service/app/api/internal/svc"
	"management-system/service/app/api/internal/types"
	"management-system/service/app/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 先根据手机号查找用户是否注册
	user, err := l.svcCtx.UsersModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "mobile:%s,err:%v", req.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("该手机号已被注册"), "", req.Mobile, err)
	}

	// 创建user数据
	users := new(models.Users)
	users.Mobile = req.Mobile
	// users.Password = req.Password
	users.Password = utils.Md5ByString(req.Password)
	users.Username = req.UserName
	users.Sex = req.Sex
	users.Email = req.Email
	users.Address = req.Address
	users.Birthday = req.Birthday

	if _, err := l.svcCtx.UsersModel.Insert(l.ctx, users); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Regiter db user Insert err:%v,user:%+v", err, user)
	}

	userInfo, err := l.svcCtx.UsersModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil {
		return nil, err
	}
	// 创建userAuth数据
	usersAuth := new(models.UsersAuth)
	usersAuth.AuthKey = req.Mobile
	usersAuth.AuthType = models.UserAuthTypeSystem
	usersAuth.UserId = userInfo.Id
	if _, err := l.svcCtx.UsersAuthModel.Insert(l.ctx, usersAuth); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Regiter db user Insert err:%v,user:%+v", err, user)
	}

	userId := usersAuth.UserId
	// 生成token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, l.svcCtx.Config.JwtAuth.AccessExpire, userId)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		AccessToken:  jwtToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil

}

func (l *RegisterLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
