package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/alicebob/miniredis"
	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

func init() {
	keeper.RedisAddr = os.Getenv("REDISADDR")
	keeper.RedisPw = os.Getenv("REDISPW")

	// Use local redis if external one does not exist
	if keeper.RedisAddr == "" {
		rds, _ := miniredis.Run()
		keeper.RedisAddr = rds.Addr()
		keeper.RedisPw = ""
	}
}

// Handler is the plugin handler
func Handler(v url.Values) (string, error) {
	text := strings.Split(v["text"][0], " ")
	return keeper.Handle(text[0], v["team_id"][0], v["user_id"][0], v["user_name"][0], text[1:])
}
