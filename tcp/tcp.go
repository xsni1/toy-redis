package tcp

import (
	"fmt"
	"net"

	"github.com/xsni1/toy-redis/parser"
)

type TCPConfig struct {
	Addr string
	Host string
}

// when reading from tcp socket
// we never know when to stop - unless client disconnects
// this is why we need some kind of protocol - rules about the shape of the data - how it starts, ends etc.
// but at the same time we have to treat tcp as a stream - so we have to keep reading it until we have our defined by the protocl end of message
// i think there should also be some kind of timeout so we aren't stuck reading forever (redis does not do it!!)
func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	// TODO: check what's redis max message
	buf := make([]byte, 4096)

	for {
		// Could very well be simplified to not use goroutines at all
		// but wanted to mess around
		// TODO: move it all to `Parse` method?
		in := make(chan []byte)
		go func() {
			for {
				n, err := conn.Read(buf)
				if err != nil {
					close(in)
					return
				}
				in <- buf[:n]
			}
		}()

		out := parser.Parse(in)
		res := <-out
		close(in)
		fmt.Println(res)
		// execute cmd
	}
}

func ListenTCPSocket(config TCPConfig) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", config.Addr, config.Host))
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
		go handleConn(conn)
	}
}
