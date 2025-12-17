package server_test

import (
	"log/slog"
	"os"
	"testing"

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

	testClientActions()
}

func testClientActions() {
	client1 := client_test.NewTestClient("client1", "auth tester then bad message sender", 0)
	client1.ConnectToServer()
	client1.TestUnauthenticated()
	client1.ConnectToServer()

	client1.TestAuthenticationFailure()
	client1.TestAuthenticationSuccess()

	client1.Close()
}