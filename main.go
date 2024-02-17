package main

import (
	"fmt"
	"net"
	"slices"
)

func encode() {

}

// when reading from tcp socket
// we never know when to stop - unless client disconnects
// this is why we need some kind of protocol - rules about the shape of the data - how it starts, ends etc.
// but at the same time we have to treat tcp as a stream - so we have to keep reading it until we have our defined end
// i think there should also be some kind of timeout so we aren't stuck reading forever (redis does not do it!!)
func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	var (
		buffer  []byte
		cur     int
		command string
	)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("err reading msg: %v", err)
			return
		}
		fmt.Printf("reading %d bytes, err: %v\n", n, err)

		// parse command
		if cur == 0 {
			// multi bulk
			if buffer[0] == '*' {
				numend := slices.Index(buffer, '\r')
				if numend == -1 || len(buffer) <= numend+1 || (len(buffer) > numend+1 && buffer[numend+1] != '\n') {
                    // i need two buffers, so in this cause i can just call continue, and the read will concat new data with old
				}
			} else {
				// inline
			}
		}
		// if readlen == 0 {
		// 	switch msgtype := buffer[0]; msgtype {
		// 	case '+':

		// 	case ':':
		// 	}
		// }
		// readlen += n
	}

	// conn.Write([]byte("asd"))
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:3456")
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
		go handleConn(conn)
	}
}
