package chat

import (
	"errors"
	"fmt"
	"github.com/bysir-zl/bygo/util/auth"
	"strconv"
)

// return a user that only set id
func VerifyUser(token string) (uid int64, err error) {
	data, errCode := auth.JWTDecode(token, SECRET)
	if errCode != 0 {
		err = fmt.Errorf("code is %d", errCode)
		return
	}
	if data.Iss != SERVER {
		err = errors.New("verify error")
		return
	}
	uid, _ = strconv.ParseInt(data.Sub, 10, 64)
	return
}
