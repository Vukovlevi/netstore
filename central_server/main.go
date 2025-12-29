package main

import (
	"log/slog"
	"os"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	server := tcp.NewServer()
	server.Start()
}
