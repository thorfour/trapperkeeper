package keeper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
)

var (
	NoCmdErr   = fmt.Errorf("command not found")
	NumArgsErr = fmt.Errorf("wrong number of args")
	NoActivWin = fmt.Errorf("no window active. Try the new command")
)

const (
	new     = "new"
	add     = "add"
	release = "release"

	timeLayout = "Jan 2 9:04PM MST" // time layout to parse

	windowKey = "window" // currently active window key

	redisAddr = ""
	redisPw   = ""
)

var lookup = map[string]func(string, string, []string) error{
	new:     newWindow,
	add:     addSubmission,
	release: releaseWindow,
}

// Handle will process a given command
func Handle(cmd, uid, uname string, args []string) error {
	// Lookup the command
	f, ok := lookup[cmd]
	if !ok {
		return NoCmdErr
	}

	// Execute the command
	return f(uid, uname, args)
}

// release window will return the submissions from the current window
func releaseWindow(uid, uname string, args []string) error {
	c := connectRedis()
	w, err := getWindow(c)
	if err != nil {
		return NoActivWin
	}

	// Ensure the window has expired
	if time.Now().Unix() < w.Expire {
		return fmt.Errorf("Window has not expired")
	}

	var s string
	for _, u := range w.Submissions {
		s = fmt.Sprintf("%s%s : %s\n", s, u.Name, u.Submission)
	}

	fmt.Println(s)
	return nil
}

// addSubmission will add a submission to the current window
func addSubmission(uid, uname string, args []string) error {
	c := connectRedis()
	w, err := getWindow(c)
	if err != nil {
		return NoActivWin
	}

	// ensure the time hasn't expired
	if time.Now().Unix() > w.Expire {
		return fmt.Errorf("submission window has expired")
	}

	s := strings.Join(args, " ") // rejoin the args into the submission string
	w.Submissions[uid] = user{uid, uname, s}

	fmt.Fprintf(os.Stderr, "Submitted: %s\n", s)

	// save the window
	return saveWindow(c, w)
}

// newWindow creates a new user submission window
func newWindow(uid, uname string, args []string) error {
	c := connectRedis()

	// Ensure the current window is not active
	w, err := getWindow(c)
	if err == nil {
		if time.Now().Unix() < w.Expire {
			return fmt.Errorf("current window is still active until %v", time.Unix(w.Expire, 0).Format(timeLayout))
		}
	}

	// first time creating a window
	if err != nil {
		w = &window{}
	}

	// parse the window duration
	duration, err := time.ParseDuration(strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("bad time format: %v", err)
	}

	// Set the new window
	w.Expire = time.Now().Add(duration).Unix()
	w.Owner = uid
	w.Submissions = make(map[string]user)

	fmt.Printf("New pick window expires at %s\n", time.Unix(w.Expire, 0).Format(timeLayout))

	// Save the window
	return saveWindow(c, w)
}

// saveWindow will save the window to redis
func saveWindow(client *redis.Client, w *window) error {
	serialized, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("unable to encode: %v", err)
	}

	_, err = client.Set(windowKey, string(serialized), 0).Result()
	return err
}

// getWindow returns the current window if one exists
func getWindow(client *redis.Client) (*window, error) {
	serialized, err := client.Get(windowKey).Result()
	if err != nil {
		return nil, err
	}

	// Unserialize the window struct
	w := &window{}
	err = json.Unmarshal([]byte(serialized), w)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %v", err)
	}

	return w, nil
}

func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPw,
		DB:       0,
	})
}
