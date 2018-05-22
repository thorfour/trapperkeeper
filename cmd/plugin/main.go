package main

import (
	"net/url"
	"strings"

	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

// Handler is the plugin handler
func Handler(v url.Values) (string, error) {
	text := strings.Split(v["text"][0], " ")
	return keeper.Handle(text[0], v["user_id"][0], v["user_name"][0], text[1:])
}
