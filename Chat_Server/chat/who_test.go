package chat

import (
	"strings"
	"testing"

	zt "zirc/tests"
)

func TestWHO_NickQuery(t *testing.T) {
	// override ports so if a local Docker instance is running
	// ports don't conflict
	t.Setenv("IRC_S2S_PORT", "7001")
	serverAddr := "127.0.0.1:6677"

	is := NewIrcServer(
		"zirc-test.com",
		"vtest",
		serverAddr,
		"leaf",
		nil,
		nil,
		nil,
	)
	go is.Run()

	// time.Sleep(5 * time.Second)

	// mc := zt.MockClient{}
	// mc2 := zt.MockClient{}
	// mc3 := zt.MockClient{}
	// mc4 := zt.MockClient{}
	// mc5 := zt.MockClient{}
	// mc6 := zt.MockClient{}
	// mc7 := zt.MockClient{}
	// _ = mc.Register(is.Addr, "zak", "zak")
	// _ = mc2.Register(is.Addr, "zik", "zik")
	// _ = mc3.Register(is.Addr, "zuk", "zuk")
	// _ = mc4.Register(is.Addr, "marmar", "marmar")
	// _ = mc5.Register(is.Addr, "sadie", "sadie")
	// _ = mc6.Register(is.Addr, "kk", "kk")
	// _ = mc7.Register(is.Addr, "zorbra", "zorbra")

	t.Run("Test 1: single NICK query with no JOINed channels", func(t *testing.T) {
		// t.Parallel()
		mc_recv := make(chan zt.MockResponse)
		mc := zt.NewMockClient(serverAddr, mc_recv)
		go mc.Run()

		mc2_recv := make(chan zt.MockResponse)
		mc2 := zt.NewMockClient(serverAddr, mc2_recv)
		go mc2.Run()

		mc_reg := zt.RegisterMsg(serverAddr, "zak", "zak")
		mc.Send(mc_reg, "004", false)
		for {
			res := <-mc_recv
			if strings.Contains(res.Response, "004") {
				break
			}
		}

		mc2_reg := zt.RegisterMsg(serverAddr, "marmar", "marmar")
		mc2.Send(mc2_reg, "004", false)
		for {
			res := <-mc_recv
			if strings.Contains(res.Response, "004") {
				break
			}
		}

		mc2.Send("WHO zak\r\n", "", false)
		_got := <-mc2_recv
		got := _got.Response
		want := ":zirc-test.com 352 marmar * zak cloak.z.irc zirc-test.com zak H :0 zak \r\n"

		if want != got {
			t.Errorf("unexpected response:\ndont_want: '%s'\ngot: '%s'", want, got)
		}
	})

	// t.Run("Test 2: single NICK with JOINed channels", func(t *testing.T) {
	// 	t.Parallel()

	// 	_ = mc2.Send(is.Addr, "JOIN #test\r\n", false)

	// 	want := ":localhost 451 * :You have not registered \r\n"
	// 	got := mc3.Send(is.Addr, "WHO zuk\r\n", true)

	// 	if want != got {
	// 		t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", want, got)
	// 	}
	// })

	// t.Run("Test 3: single NICK with regex ?", func(t *testing.T) {
	// 	t.Parallel()

	// 	want := ":localhost 451 * :You have not registered \r\n"
	// 	got := mc5.Send(is.Addr, "WHO z?k\r\n", true)

	// 	if want != got {
	// 		t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", want, got)
	// 	}
	// })

	// t.Run("Test 4: NICK with regex *", func(t *testing.T) {
	// 	t.Parallel()

	// 	want := ":localhost 451 * :You have not registered \r\n"
	// 	got := mc6.Send(is.Addr, "WHO z*\r\n", true)

	// 	if want != got {
	// 		t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", want, got)
	// 	}
	// })

	// t.Run("Test 5: NICK with regex *", func(t *testing.T) {
	// 	t.Parallel()

	// 	want := `
	// 	:zirc-test.com 352 kk * zach cloak.z.irc zirc-test.com zach H :0\r\n
	// 	:zirc-test.com 352 kk * zak cloak.z.irc zirc-test.com zak H :0\r\n
	// 	:zirc-test.com 352 kk * zik cloak.z.irc zirc-test.com zik H :0\r\n
	// 	:zirc-test.com 352 kk * zuk cloak.z.irc zirc-test.com zuk H :0\r\n
	// 	:zirc-test.com 352 kk * dd cloak.z.irc zirc-test.com zorbra H :0\r\n
	// 	`
	// 	got := mc6.Send(is.Addr, "WHO z*\r\n", true)

	// 	if want != got {
	// 		t.Errorf("unexpected response:\ndo_not_want: '%s'\ngot: '%s'", want, got)
	// 	}
	// })
}
