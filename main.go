package main

import (
	"github.com/xsni1/toy-redis/store"
	"github.com/xsni1/toy-redis/tcp"
)

func main() {
	// TODO: provide config by file
	tcpConfig := tcp.Config{
		Addr: "127.0.0.1",
		Host: "6379",
	}
	storeConfig := store.Config{
		Appendonly: true,
	}
	store := store.NewStore(storeConfig)
	srv := tcp.NewServer(tcpConfig, &store)

	srv.ListenTCPSocket(tcpConfig)
}
