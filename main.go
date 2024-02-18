package main

import (
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"sync"
)

// map could be sharded to minimize time spent waiting for locks
// or not? it could be sharded if we were to use simple map with hand-made locks
var store = sync.Map{}

func parseCommand(elements []string) (any, error) {
	switch elements[0] {
	case "SET":
		store.Store(elements[1], elements[2])
		return "OK", nil
	case "GET":
		if val, ok := store.Load(elements[1]); ok {
			return val, nil
		}
		return "", fmt.Errorf("not found")
	}
	return "", fmt.Errorf("failure during command parsing")
}

// when reading from tcp socket
// we never know when to stop - unless client disconnects
// this is why we need some kind of protocol - rules about the shape of the data - how it starts, ends etc.
// but at the same time we have to treat tcp as a stream - so we have to keep reading it until we have our defined end
// i think there should also be some kind of timeout so we aren't stuck reading forever (redis does not do it!!)
func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	var (
		// co jesli wiadomosc jest wieksza? ucinane czy porcjowane
		tempBuffer   = make([]byte, 1024)
		buffer       []byte
		elements     []string
		multiBulkLen int
	)

	for {
		n, err := conn.Read(tempBuffer)
		buffer = append(buffer, tempBuffer[:n]...)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("disconnecting client\n")
				return
			}
			fmt.Printf("err reading msg: %v", err)
			return
		}

		if multiBulkLen == 0 {
			// multi bulk
			if buffer[0] == '*' {
				numend := slices.Index(buffer, '\r')
				if numend == -1 || len(buffer) <= numend+1 || (len(buffer) > numend+1 && buffer[numend+1] != '\n') {
					continue
				}
				multiBulkLen, err = strconv.Atoi(string(buffer[1:numend]))
				if err != nil {
					fmt.Printf("err reading array len: %v\n", err)
					return
				}
				// 2 to put pointer on the first byte after \n
				buffer = buffer[numend+2:]
			} else {
				// inline
			}
		}

		for multiBulkLen > 0 {
			fmt.Println(buffer)
			if buffer[0] != '$' {
				fmt.Print("err decoding, expected $\n")
				return
			}
			numend := slices.Index(buffer, '\r')
			if numend == -1 || len(buffer) <= numend+1 || (len(buffer) > numend+1 && buffer[numend+1] != '\n') {
				// i need two buffers, so in this cause i can just call continue, and the read will concat new data with old
				break
			}
			strLen, err := strconv.Atoi(string(buffer[1:numend]))
			if err != nil {
				fmt.Printf("err reading element len: %v\n", err)
				return
			}
			if len(buffer) <= numend+2+strLen {
				break
			}
			elements = append(elements, string(buffer[numend+2:numend+2+strLen]))
			buffer = buffer[numend+2+strLen+2:]
			multiBulkLen--
			fmt.Println("parsed", elements)
		}

		if multiBulkLen > 0 {
			continue
		}

		fmt.Println("parsed: ", elements)
		// parse command
		res, err := parseCommand(elements)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
		} else {
			conn.Write([]byte(fmt.Sprintf("+%s\r\n", res)))
		}

		buffer = []byte{}
		elements = []string{}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
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
