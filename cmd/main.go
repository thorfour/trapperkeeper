package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thorfour/trapperkeeper/pkg/keeper"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "expected arg")
		return
	}

	m := make(map[string]string)
	if err := json.Unmarshal([]byte(os.Args[1]), &m); err != nil {
		fmt.Fprint(os.Stderr, "json.Unmarshal: ", err)
		return
	}

	text := strings.Split(m["text"], " ")

	if err := keeper.Handle(text[0], m["user_id"], m["user_name"], text[1:]); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
}
