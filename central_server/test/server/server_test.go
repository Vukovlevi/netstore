package server_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/vukovlevi/netstore/central_server/tcp"
	client_test "github.com/vukovlevi/netstore/central_server/test/client"
)

func TestServer(t *testing.T) {
	logFile, err := os.Create("main.log")
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})))

	server := tcp.NewServer()
	go server.Start()

	testClientAuthActions()
    testClientActions()
}

func testClientAuthActions() {
	client1 := client_test.NewTestClient("client1", "auth tester then bad message sender", 0)
	client1.ConnectToServer()
	client1.TestUnauthenticated()
	client1.ConnectToServer()

	client1.TestAuthenticationFailure()
	client1.TestAuthenticationSuccess()

	client1.Close()
}

func testClientActions() {
	client1 := client_test.NewTestClient("client1", "auth tester then bad message sender", 0)
	client2 := client_test.NewTestClient("client2", "test client number 2", 5)
	client3 := client_test.NewTestClient("client3", "test client number 3", 12)

    client1.ConnectToServer()
    client2.ConnectToServer()
    client3.ConnectToServer()
    client1.TestAuthenticationSuccess()
    client2.TestAuthenticationSuccess()
    client3.TestAuthenticationSuccess()

    client1.SendMalformedMessage()
    client2.SendSearchRequest()

    client3.ReadClientSearch()
    client3.SendAnswer(time.Microsecond)

    client2.ReadClientAnswer(false, "client3", "test client number 3", 13)

    client1.SendMalformedMessage()
    client3.SendSearchRequest()

    client2.ReadClientSearch()
    client2.SendAnswer(time.Second * 4)

    client3.ReadClientAnswer(true, "client2", "test client number 2", 6)

    client1.Close()
    client2.Close()
    client3.Close()
}
