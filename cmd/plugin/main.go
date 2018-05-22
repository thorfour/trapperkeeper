package main

import (
	"net/url"
	"strings"

	"github.com/alicebob/miniredis"
	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

var (
	setupErr error
)

func init() {
	var rds *miniredis.Miniredis
	rds, setupErr = miniredis.Run()
	keeper.RedisAddr = rds.Addr()
}

// Handler is the plugin handler
func Handler(v url.Values) (string, error) {
	if setupErr != nil {
		return "", setupErr
	}
	text := strings.Split(v["text"][0], " ")
	return keeper.Handle(text[0], v["user_id"][0], v["user_name"][0], text[1:])
}
