package keeper

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

func TestMain(m *testing.M) {
	rds, err := miniredis.Run()
	if err != nil {
		fmt.Println("Failed to start miniredis")
		os.Exit(1)
	}
	defer rds.Close()
	RedisAddr = rds.Addr()
	os.Exit(m.Run())
}

func TestHumanReadableTime(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split("3h", " ")
	fmt.Println("New Window: ", expire)

	if _, err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}
}

func TestNoWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	_, err := Handle(add, "0", "test", []string{"submission0"})
	if err != errNoActivWin {
		t.Error("wrong error returned when one was expected: ", err)
	}
}

func TestAddNewWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split("1h", " ")
	fmt.Println("New Window: ", expire)

	if _, err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	_, err := Handle(add, "0", "test", []string{"submission2"})
	if err != nil {
		t.Error("failed to add to newly created window: ", err)
	}
}

func TestAddExpiredWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split("1s", " ")
	fmt.Println("New Window: ", expire)

	if _, err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	time.Sleep(2 * time.Second)

	_, err := Handle(add, "0", "test", []string{"submission2"})
	if err == nil {
		t.Error("added to expired window")
	}
}

func TestReleaseExpiredWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split("3s", " ")
	fmt.Println("New Window: ", expire)

	if _, err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	for i := 0; i < 10; i++ {
		_, err := Handle(add, fmt.Sprintf("%v", i), "test", []string{fmt.Sprintf("submit%v", i)})
		if err != nil {
			t.Error("error added to valid window: ", err)
		}
	}

	time.Sleep(3 * time.Second)

	if _, err := Handle(release, "0", "test", []string{}); err != nil {
		t.Error("failed to release expired window")
	}
}

func TestReleaseNotExpired(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split("3h", " ")
	fmt.Println("New Window: ", expire)

	if _, err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	for i := 0; i < 10; i++ {
		_, err := Handle(add, fmt.Sprintf("%v", i), "test", []string{fmt.Sprintf("submit%v", i)})
		if err != nil {
			t.Error("error added to valid window: ", err)
		}
	}

	time.Sleep(3 * time.Second)

	if _, err := Handle(release, "0", "test", []string{}); err == nil {
		t.Error("release unexpired window")
	}
}

func TestReleaseInvalid(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	if _, err := Handle(release, "0", "test", []string{}); err == nil {
		t.Error("expected error when releasing before creating a window")
	}
}
