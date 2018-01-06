package keeper

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNoWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	err := Handle(add, "0", "test", []string{"submission0"})
	if err != NoActivWin {
		t.Error("wrong error returned when one was expected: ", err)
	}
}

func TestAddNewWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split(time.Now().Add(100*time.Minute).Format(timeLayout), " ")
	fmt.Println("New Window: ", expire)

	if err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	err := Handle(add, "0", "test", []string{"submission2"})
	if err != nil {
		t.Error("failed to add to newly created window: ", err)
	}
}

func TestAddExpiredWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split(time.Now().Add(1*time.Second).Format(timeLayout), " ")
	fmt.Println("New Window: ", expire)

	if err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	time.Sleep(2 * time.Second)

	err := Handle(add, "0", "test", []string{"submission2"})
	if err == nil {
		t.Error("added to expired window")
	}
}

func TestReleaseExpiredWindow(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split(time.Now().Add(3*time.Second).Format(timeLayout), " ")
	fmt.Println("New Window: ", expire)

	if err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	for i := 0; i < 10; i++ {
		err := Handle(add, fmt.Sprintf("%s", i), "test", []string{fmt.Sprintf("submit%v", i)})
		if err != nil {
			t.Error("error added to valid window: ", err)
		}
	}

	time.Sleep(3 * time.Second)

	if err := Handle(release, "0", "test", []string{}); err != nil {
		t.Error("failed to release expired window")
	}
}

func TestReleaseNotExpired(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	expire := strings.Split(time.Now().Add(3*time.Hour).Format(timeLayout), " ")
	fmt.Println("New Window: ", expire)

	if err := Handle(new, "0", "test", expire); err != nil {
		t.Error("failed to add new window: ", err)
	}

	for i := 0; i < 10; i++ {
		err := Handle(add, fmt.Sprintf("%s", i), "test", []string{fmt.Sprintf("submit%v", i)})
		if err != nil {
			t.Error("error added to valid window: ", err)
		}
	}

	time.Sleep(3 * time.Second)

	if err := Handle(release, "0", "test", []string{}); err == nil {
		t.Error("release unexpired window")
	}
}

func TestReleaseInvalid(t *testing.T) {
	c := connectRedis()
	c.Del(windowKey)

	if err := Handle(release, "0", "test", []string{}); err == nil {
		t.Error("expected error when releasing before creating a window")
	}
}
