package qiniu

import (
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func GetToken() string {
	Load()
	var maxInt uint64 = 1 << 32
	putPolicy := storage.PutPolicy{
		Scope:   Q.Bucket,
		Expires: maxInt,
	}
	fmt.Println(Q.Bucket,Q.AccessKey,Q.SecretKey)
	mac := qbox.NewMac(Q.AccessKey, Q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
