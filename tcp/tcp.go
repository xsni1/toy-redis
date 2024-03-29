package tcp

import (
	"fmt"
	"net"

	"github.com/xsni1/toy-redis/core"
	"github.com/xsni1/toy-redis/parser"
	"github.com/xsni1/toy-redis/store"
)

type Config struct {
	Addr string
	Port string
}

type Server struct {
	config Config
	store  *store.Store
	core   *core.Core
}

func NewServer(config Config, store *store.Store, core *core.Core) Server {
	return Server{
		config: config,
		store:  store,
		core:   core,
	}
}

// when reading from tcp socket
// we never know when to stop - unless client disconnects
// this is why we need some kind of protocol - rules about the shape of the data - how it starts, ends etc.
// but at the same time we have to treat tcp as a stream - so we have to keep reading it until we have our defined by the protocl end of message
// i think there should also be some kind of timeout so we aren't stuck reading forever (redis does not do it!!)
func (s *Server) handleConn(conn *net.TCPConn) {
	defer conn.Close()
	// Could very well be simplified to not use goroutines at all
	// but wanted to mess around
	// TODO: move it all to `Parse` method?
	// TODO: check what's redis max message

	in := make(chan []byte)
	buf := make([]byte, 4096)
	done := make(chan bool)
	go func() {
		for {
			n, err := conn.Read(buf)
			if err != nil {
				close(in)
				done <- true
				break
			}
			in <- buf[:n]
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			out := parser.Parse(in)
			res := <-out
			conn.Write(s.core.Execute(res))
		}
	}
}

func (s *Server) ListenTCPSocket(config Config) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", config.Addr, config.Port))
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("err accepting: %v", err)
			continue
		}
		// TODO: Alternate single-threaded version using epoll to mirror original redis implementation
		go s.handleConn(conn)
	}
}
