package main

import (
	"github.com/xsni1/toy-redis/tcp"
)

func main() {
	// TODO: provide config by file
	config := tcp.TCPConfig{
		Addr: "127.0.0.1",
		Host: "6379",
	}

	tcp.ListenTCPSocket(config)
}
