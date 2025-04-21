package chat

import (
	"fmt"
	"os"
	"testing"
	"zirc/tests"
)

func TestPASS_Basic(t *testing.T) {
	val := os.Getenv("CHAT_ROOT")
	fmt.Println(val)
	t.Setenv("CONFIG_FILE", "../config.yaml")

	cmd_param_slice := map[string]interface{}{}
	mockConn := tests.MockConn{}
	client, _, _ := NewClient("test", "test", nil, mockConn, false)
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

	t.Run("Test 3: Password already sent & accepted by client", func(t *testing.T) {
		cmd_param_slice["client"] = client
		cmd_param_slice["params"] = ""
		client.session.state["server_password"] = "accepted"

		got := cmd.Fn(cmd_param_slice)
		want := ":You may not reregister"
		wantCode := "462"

		if want != got.Msg() {
			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
		}

		if wantCode != got.Code() {
			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
		}

		delete(client.session.state, "server_password")
	})
}
