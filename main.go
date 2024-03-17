package main

import (
	"github.com/xsni1/toy-redis/core"
	"github.com/xsni1/toy-redis/store"
	"github.com/xsni1/toy-redis/tcp"
)

// TODO: provide config by file
type appConfig struct {
	tcp struct {
		addr string
		port string
	}
	store struct {
		appendonly bool
	}
}

var config = appConfig{
	tcp: struct {
		addr string
		port string
	}{
		addr: "127.0.0.1",
		port: "6379",
	},
	store: struct {
		appendonly bool
	}{
		appendonly: true,
	},
}

func main() {
	tcpConfig := tcp.Config{
		Addr: config.tcp.addr,
		Port: config.tcp.port,
	}
	storeConfig := store.Config{
		Appendonly: config.store.appendonly,
	}
	store := store.NewStore(storeConfig)
	core := core.NewCore(store)
	srv := tcp.NewServer(tcpConfig, store, core)

	srv.ListenTCPSocket(tcpConfig)
}
