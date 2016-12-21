package main

import (
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/chat"
)

func main() {
	s := chat.Server{}
	log.Info("chat", "start success")
	err := s.Server()
	panic(err)
}
