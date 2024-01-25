package qiniu

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniu struct {
	AccessKey string
	SecretKey string
	Bucket    string // 存储空间
	Domain    string // 域名
}

var Q Qiniu

func Load() {
	Q = Qiniu{
		AccessKey: "",
		SecretKey: "",
		Bucket:    "",
		Domain:    "",
	}
}

func UploadQiniu(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()

	putPolicy := storage.PutPolicy{
		Scope: Q.Bucket,
	}

	mac := qbox.NewMac(Q.AccessKey, Q.SecretKey)

	// 获取上传凭证
	upToken := putPolicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          
		UseCdnDomains: false,               // 是否使用CDN上传加速
		UseHTTPS:      false,               // 是否使用https
	}

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传返回后的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 自定义文件名及后缀
	key := "(" + time.Now().String() + ")" + file.Filename

	if err := formUploader.Put(context.Background(), &ret,
		upToken, key, src, file.Size, &putExtra); err != nil {
		return "", err
	}

	return Q.Domain + "/" + ret.Key, nil
}
