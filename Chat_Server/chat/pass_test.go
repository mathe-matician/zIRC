package chat

import (
	"testing"
	"zirc/tests"
)

func TestPASS_Basic(t *testing.T) {
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
			t.Errorf("unexpected code:\nwant: '%s'\ngot: '%s'", wantCode, got)
		}

		delete(client.session.state, "server_password")
	})
}

func TestPASS_Server(t *testing.T) {
	t.Setenv("IRC_S2S_CONFIG_FILE", "")
	// TODO
	// servermanagerconfig throws an error here because it is trying to read the s2s config even though it _should_ be empty
	is := NewIrcServer(
		"zirc-test.com",
		"vtest",
		"127.0.0.1:6677",
		"leaf",
		nil,
		nil,
		nil,
	)
	g_Server = is

	cmd_param_slice := map[string]interface{}{}
	mockConn := tests.MockConn{}
	client, _, _ := NewClient("test", "test", nil, mockConn, false)
	client.IsServer = true
	cmd_param_slice["client"] = client
	cmd_param_slice["params"] = ""

	cmd := command_map["PASS"]

	t.Run("Test 1: terminate not whitelisted server", func(t *testing.T) {
		got := cmd.Fn(cmd_param_slice)
		want := "TERMINATE"
		wantCode := "666"

		if want != got.Msg() {
			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
		}

		if wantCode != got.Code() {
			t.Errorf("unexpected code:\nwant: '%s'\ngot: '%s'", wantCode, got)
		}
	})
}
