package commands_tests

// import (
// 	"testing"

// 	zc "zirc/chat"
// 	zt "zirc/tests"
// )

// func TestPASS_Disabled(t *testing.T) {
// 	is := zc.NewIrcServer(
// 		"zirc-test.com",
// 		"vtest",
// 		"127.0.0.1:6677",
// 		"leaf",
// 		nil,
// 		nil,
// 		nil,
// 	)
// 	go is.Run()

// 	t.Run("Test 1: dont get the pw required response", func(t *testing.T) {
// 		t.Parallel()
// 		mc := zt.MockClient{}
// 		dont_want := ":localhost 451 * :You need to send your password before registering \r\n"
// 		got := mc.Send(is.Addr, "JOIN \r\n", false)

// 		if dont_want == got {
// 			t.Errorf("unexpected response:\ndont_want: '%s'\ngot: '%s'", dont_want, got)
// 		}
// 	})

// 	// NOTE - this test should timeout as no response should ever be sent back from the server
// 	t.Run("Test 2: when no pw enabled nothing returned when PASS run", func(t *testing.T) {
// 		t.Parallel()
// 		mc := zt.MockClient{}
// 		dont_want := ":localhost 451 * :You have not registered \r\n"
// 		got := mc.Send(is.Addr, "PASS \r\n", true)

// 		if dont_want == got {
// 			t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", dont_want, got)
// 		}
// 	})
// }

// func TestPASS_Enabled(t *testing.T) {
// 	t.Setenv("IRC_SERVER_PASSWORD", "$2a$10$RFRaJiN40EzL9.nvfq4uy.w6UHs7WuF3EkBpxj.xUPjNd8tUs7sfi")
// 	is := zc.NewIrcServer(
// 		"zirc-test.com",
// 		"vtest",
// 		"0.0.0.0:6677",
// 		"leaf",
// 		nil,
// 		nil,
// 		nil,
// 	)
// 	go is.Run()

// 	t.Run("Test 1: Running any cmd requires registration", func(t *testing.T) {
// 		t.Parallel()
// 		mc := zt.MockClient{}
// 		want := ":localhost 464 * :You need to send your password before registering \r\n"
// 		got := mc.Send(is.Addr, "1235 \r\n", false)

// 		if want != got {
// 			t.Errorf("unexpected response:\nwant: '%s'\ngot: '%s'", want, got)
// 		}
// 	})

// 	t.Run("Test 2: Successful pw returns nothing", func(t *testing.T) {
// 		t.Parallel()
// 		mc := zt.MockClient{}
// 		want := ""
// 		got := mc.Send(is.Addr, "PASS 1235 \r\n", true)

// 		if want != got {
// 			t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", want, got)
// 		}
// 	})
// }
