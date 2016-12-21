package chat

import (
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util/auth"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	token := auth.JWTEncode(time.Now().AddDate(1,0,0).Unix(), SERVER, "U", "2", "chat", SECRET)
	log.Info("test", token)
	token2:= auth.JWTEncode(time.Now().AddDate(1,0,0).Unix(), SERVER, "U", "3", "chat", SECRET)
	log.Info("test", token2)
}
