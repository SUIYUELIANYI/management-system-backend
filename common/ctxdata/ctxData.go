package ctxdata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// 存放解析Jwttoken获得userId，方便每个接口调用
var CtxKeyJwtUserId = "userId"

func GetUidFromCtx(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
			fmt.Println(uid)
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err: %+v", err)
		}
	}
	return uid
}
