package chat

import (
	"fmt"
	"os"
	"testing"
)

func TestPASS_Basic(t *testing.T) {
	val := os.Getenv("CHAT_ROOT")
	fmt.Println(val)
	t.Setenv("CONFIG_FILE", "../config.yaml")

	cmd_param_slice := map[string]interface{}{}
	client, _, _ := NewClient("test", "test", nil, nil, false)
	cmd := command_map["PASS"]

	t.Run("Test 1: no internal client struct passed to PASS", func(t *testing.T) {
		got := cmd.Fn(cmd_param_slice)
		want := ":Unknown error occurred"

		if want != got.Msg() {
			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
		}
	})

	t.Run("Test 2: No internal params passed to PASS", func(t *testing.T) {
		cmd_param_slice["client"] = client

		got := cmd.Fn(cmd_param_slice)
		want := ":Need more params"

		if want != got.Msg() {
			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
		}
	})
}
