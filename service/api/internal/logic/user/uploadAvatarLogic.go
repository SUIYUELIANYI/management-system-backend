package user

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"management_system/common/github_image"
	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadAvatarLogic) UploadAvatar() (resp *types.UploadAvatarResp, err error) {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	_, fileHeader, err := l.r.FormFile("file")
	if err != nil {
		return nil, response.Error(400, "文件上传错误")
	}

	PATH := "avatar"

	filepath := "./"
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	fileExt := path.Ext(filepath + fileHeader.Filename)

	fileHeader.Filename = strconv.FormatInt(userId, 10) + "_" + GetRandomString(16) + fileExt

	filename := filepath + fileHeader.Filename

	if err := SaveUploadedFile(fileHeader, filename); err != nil {
		return nil, response.Error(401, "上传图片失败")
	}

	// 上传头像
	Base64 := ImagesToBase64(filename)
	picUrl, _, _ := github_image.RepoCreate().Push(PATH, fileHeader.Filename, Base64)
	os.Remove(filename)

	url := picUrl
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, response.Error(401, err.Error())
	}
	userInfo.Avatar = url
	err = l.svcCtx.UserModel.Update(l.ctx, userInfo)
	if err != nil {
		return nil, response.Error(401, err.Error())
	}
	if url == "" {
		return nil, response.Error(401, "上传失败")
	}
	return &types.UploadAvatarResp{
		Url: url,
	}, nil
}

func ImagesToBase64(str_images string) string {
	f, _ := os.Open(str_images)
	defer f.Close()
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}

func GetRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
