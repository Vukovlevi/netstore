package main

import "github.com/vukovlevi/netstore/central_server/tcp"

func main() {
	server := tcp.NewServer()
	server.Start()
}
