package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(text string, salt string) string {
	ctx := md5.New()
	ctx.Write([]byte(salt + text))
	return hex.EncodeToString(ctx.Sum(nil))
}
