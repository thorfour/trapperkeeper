package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

func main() {
	m := make(map[string]string)
	if err := json.Unmarshal([]byte(os.Args[1]), &m); err != nil {
		fmt.Fprint(os.Stderr, "json.Unmarshal: ", err)
	}

	text := strings.Split(m["text"], " ")

	if err := keeper.Handle(text[0], m["user_id"], m["user_name"], text[1:]); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
}
