package keeper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
)

var (
	errNoCmd      = fmt.Errorf("command not found")
	errNoActivWin = fmt.Errorf("no window active. Try the new command")
	// RedisAddr is the redis enpoint to use
	RedisAddr string
	// RedisPw is if the endpoint has a password
	RedisPw string

	// windowKey that is used to lookup a window. Is set every invocation of Handle
	windowKey string
)

const (
	new     = "new"
	add     = "add"
	release = "release"
	current = "window"
	reset   = "reset"

	timeLayout = "Jan 2 3:04PM MST" // time layout to parse

	namespace = "keeper" // namespace for window keys
)

var lookup = map[string]func(string, string, []string) (string, error){
	new:     newWindow,
	add:     addSubmission,
	release: releaseWindow,
	current: currentWindow,
	reset:   resetWindow,
}

// Handle will process a given command
func Handle(cmd, teamID, uid, uname string, args []string) (string, error) {
	// Lookup the command
	f, ok := lookup[cmd]
	if !ok {
		return "", errNoCmd
	}

	if teamID == "" {
		return "", errNoCmd
	}

	windowKey = fmt.Sprintf("%v:%v", namespace, teamID)

	// Execute the command
	return f(uid, uname, args)
}

// currentWindow will print the current window
func currentWindow(uid, uname string, args []string) (string, error) {
	c := connectRedis()
	w, err := getWindow(c)
	if err != nil {
		return "", errNoActivWin
	}

	return fmt.Sprintf("Current window is active until %v\n", time.Unix(w.Expire, 0).Format(timeLayout)), nil
}

// release window will return the submissions from the current window
func releaseWindow(uid, uname string, args []string) (string, error) {
	c := connectRedis()
	w, err := getWindow(c)
	if err != nil {
		return "", errNoActivWin
	}

	// Ensure the window has expired
	if time.Now().Unix() < w.Expire {
		return "", fmt.Errorf("Window has not expired")
	}

	var s string
	for _, u := range w.Submissions {
		s = fmt.Sprintf("%s%s : %s\n", s, u.Name, u.Submission)
	}

	return s, nil
}

// addSubmission will add a submission to the current window
func addSubmission(uid, uname string, args []string) (string, error) {
	c := connectRedis()
	w, err := getWindow(c)
	if err != nil {
		return "", errNoActivWin
	}

	// ensure the time hasn't expired
	if time.Now().Unix() > w.Expire {
		return "", fmt.Errorf("submission window has expired")
	}

	s := strings.Join(args, " ") // rejoin the args into the submission string
	w.Submissions[uid] = user{uid, uname, s}

	// save the window
	if err := saveWindow(c, w); err != nil {
		return "", err
	}

	return "", fmt.Errorf("Submitted: %s\n", s)
}

// newWindow creates a new user submission window
func newWindow(uid, uname string, args []string) (string, error) {
	c := connectRedis()

	// Ensure the current window is not active
	w, err := getWindow(c)
	if err == nil {
		if time.Now().Unix() < w.Expire {
			return "", fmt.Errorf("current window is still active until %v", time.Unix(w.Expire, 0).Format(timeLayout))
		}
	}

	// first time creating a window
	if err != nil {
		w = &window{}
	}

	// parse the window duration
	duration, err := time.ParseDuration(strings.Join(args, " "))
	if err != nil {
		return "", fmt.Errorf("bad time format: %v", err)
	}

	// Set the new window
	w.Expire = time.Now().Add(duration).Unix()
	w.Owner = uid
	w.Submissions = make(map[string]user)

	resp := fmt.Sprintf("New pick window expires at %s\n", time.Unix(w.Expire, 0).Format(timeLayout))

	// Save the window
	return resp, saveWindow(c, w)
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

// deleteWindow deletes the current window
func deleteWindow(client *redis.Client) error {
	_, err := client.Del(windowKey).Result()
	return err
}

func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       0,
	})
}

// resetWindow deletes the current window and all submissions without returning any of them
func resetWindow(_, uname string, _ []string) (string, error) {
	c := connectRedis()
	err := deleteWindow(c)
	if err != nil {
		return "", fmt.Errorf("failed to delete window")
	}

	return fmt.Sprintf("%s has reset the active window", uname), nil
}
