package main

import (
	"fmt"
	"net"
	"slices"
	"strconv"
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
		// co jesli wiadomosc jest wieksza? ucinane czy porcjowane
		tempBuffer   = make([]byte, 1024)
		buffer       []byte
		elements     []string
		multiBulkLen int
	)

	for {
		n, err := conn.Read(tempBuffer)
		buffer = append(buffer, tempBuffer...)
		if err != nil {
			fmt.Printf("err reading msg: %v", err)
			return
		}
		fmt.Printf("reading %d bytes, err: %v\n", n, err)

        // fmt.Println("buffer start: ", string(buffer))

		// parse command
		if multiBulkLen == 0 {
			// multi bulk
			if buffer[0] == '*' {
				numend := slices.Index(buffer, '\r')
				if numend == -1 || len(buffer) <= numend+1 || (len(buffer) > numend+1 && buffer[numend+1] != '\n') {
					// i need two buffers, so in this cause i can just call continue, and the read will concat new data with old
					continue
				}
				// num of elements
				multiBulkLen, err = strconv.Atoi(string(buffer[1:numend]))
				if err != nil {
					fmt.Printf("err reading array len: %v\n", err)
					return
				}
				fmt.Printf("multibulklen: %d\n", multiBulkLen)
				// 2 to put pointer on the first byte after \n
                // fmt.Println("buffer before: ", buffer)
				buffer = buffer[numend+2:]
				// cur = numend + 2
			} else {
				// inline
			}
		}

		// fmt.Println("buffer after multibulklen: ", (buffer))

		for multiBulkLen > 0 {
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
            fmt.Println("for", multiBulkLen, buffer)
			multiBulkLen--
		}
        
        if multiBulkLen > 0 {
            continue
        }

		fmt.Println("parsed: ", elements)
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
