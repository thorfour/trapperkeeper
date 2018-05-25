package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestSimpleQuote(t *testing.T) {

	form := url.Values{
		"command":   {"/pick"},
		"text":      {"window"},
		"user_name": {"thor"},
		"user_id":   {"thor"},
	}

	body := bytes.NewBufferString(form.Encode())
	resp, err := http.Post("http://localhost:8088/v1", "application/x-www-form-urlencoded", body)
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read resp body: %v", err)
	}

	r := new(response)
	if err := json.Unmarshal(b, r); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	fmt.Println(r)
}
