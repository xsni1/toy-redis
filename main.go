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

// redis stores all the commands in json files
func parseCommand(elements []string) (string, error) {
	switch elements[0] {
	case "SET":
		store.Store(elements[1], elements[2])
		return "+OK\r\n", nil
	case "GET":
		if val, ok := store.Load(elements[1]); ok {
			return fmt.Sprintf("+%s\r\n", val), nil
		}
		return "", fmt.Errorf("not found")
	case "COMMAND":
		if elements[1] == "DOCS" {
			return "+OK\r\n", nil
		}
		return "", fmt.Errorf("error parsing")
	case "EXISTS":
		var res int
		for _, v := range elements[1:] {
			if _, b := store.Load(v); b {
				res++
			}
		}
		return fmt.Sprintf(":%d\r\n", res), nil
	case "DEL":
		var res int
		for _, v := range elements[1:] {
			if _, b := store.LoadAndDelete(v); b {
				res++
			}
		}
		return fmt.Sprintf(":%d\r\n", res), nil
	case "INCR":

	}
	// conn.Write([]byte(fmt.Sprintf("+%s\r\n", res)))
	return "", fmt.Errorf("failure during command parsing")
}

// when reading from tcp socket
// we never know when to stop - unless client disconnects
// this is why we need some kind of protocol - rules about the shape of the data - how it starts, ends etc.
// but at the same time we have to treat tcp as a stream - so we have to keep reading it until we have our defined by the protocl end of message
// i think there should also be some kind of timeout so we aren't stuck reading forever (redis does not do it!!)
func handleConn(conn *net.TCPConn) {
	defer conn.Close()
	var (
		tempBuffer   = make([]byte, 1024)
		buffer       = []byte{}
		elements     []string
		multiBulkLen int
	)

	for {
		n, err := conn.Read(tempBuffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("disconnecting client...\n")
				return
			}
			fmt.Printf("err reading msg: %v\ndisconnection client...\n", err)
			return
		}
		buffer = append(buffer, tempBuffer[:n]...)

		if multiBulkLen == 0 {
			// Multi bulk - array of bulk strings
			if buffer[0] == '*' {
				newlineIdx := slices.Index(buffer, '\r')
				if newlineIdx == -1 || len(buffer) <= newlineIdx+1 || (len(buffer) > newlineIdx+1 && buffer[newlineIdx+1] != '\n') {
					continue
				}
				multiBulkLen, err = strconv.Atoi(string(buffer[1:newlineIdx]))
				if err != nil {
					fmt.Printf("err reading array len: %v\n", err)
					return
				}
				buffer = buffer[newlineIdx+2:]
			} else {
				// inline
			}
		}

		for multiBulkLen > 0 {
			if buffer[0] != '$' {
				fmt.Print("err decoding, expected $\n")
				return
			}
			newlineIdx := slices.Index(buffer, '\r')
			if newlineIdx == -1 || len(buffer) <= newlineIdx+1 || (len(buffer) > newlineIdx+1 && buffer[newlineIdx+1] != '\n') {
				break
			}
			strLen, err := strconv.Atoi(string(buffer[1:newlineIdx]))
			if err != nil {
				fmt.Printf("err reading element len: %v\n", err)
				return
			}
			if len(buffer) <= newlineIdx+2+strLen {
				break
			}
			elements = append(elements, string(buffer[newlineIdx+2:newlineIdx+2+strLen]))
			buffer = buffer[newlineIdx+2+strLen+2:]
			multiBulkLen--
		}

		if multiBulkLen > 0 {
			continue
		}

		res, err := parseCommand(elements)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
		} else {
			conn.Write([]byte(res))
		}

		// could probably reuse existing slices
		buffer = []byte{}
		elements = []string{}
	}
}

func main() {
	// f, _ := os.Create("toy-redis.prof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

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
		// TODO: Alternate single-threaded version using epoll to mirror original redis implementation
		go handleConn(conn)
	}
}
