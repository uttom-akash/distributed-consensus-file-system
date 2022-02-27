package secsuit

import (
	"crypto/md5"
	"encoding/hex"
)

func ComputeHash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}
