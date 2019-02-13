package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thorfour/sillyputty/pkg/sillyputty"
	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
	dataDir   = "."
)

var (
	port      = flag.Int("p", 443, "port to serve on")
	redisAddr string
	redisPw   string
)

func init() {
	redisAddr = os.Getenv("REDISADDR")
	redisPw = os.Getenv("REDISPW")
	flag.Parse()
}

func main() {
	logrus.Info("Starting trapperkeeper server")

	keeper.RedisAddr = redisAddr
	keeper.RedisPw = redisPw

	s := sillyputty.New("/v1",
		sillyputty.HandlerOpt("/pick", func(v url.Values) (string, error) {
			if v == nil {
				return "", fmt.Errorf("not enough arguments")
			}

			text := strings.Split(v["text"][0], " ")
			return keeper.Handle(text[0], v["team_id"][0], v["user_id"][0], v["user_name"][0], text[1:])
		}),
	)

	s.Port = *port
	s.Run()
}
