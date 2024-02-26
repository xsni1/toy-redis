package parser

import (
	"fmt"
	"slices"
	"strconv"
)

func Parse(in <-chan []byte) <-chan []string {
	out := make(chan []string)

	go func() {
        defer close(out)
		var (
			buffer       = []byte{}
			multiBulkLen = 0
			elements     []string
		)

		for msg := range in {
			fmt.Println("Processing: ", msg)
			buffer = append(buffer, msg...)

			if multiBulkLen == 0 {
				// Multi bulk - array of bulk strings
				if buffer[0] == '*' {
					newlineIdx := slices.Index(buffer, '\r')
					if newlineIdx == -1 || len(buffer) <= newlineIdx+1 || (len(buffer) > newlineIdx+1 && buffer[newlineIdx+1] != '\n') {
						continue
					}
					c, err := strconv.Atoi(string(buffer[1:newlineIdx]))
					if err != nil {
						fmt.Printf("err reading array len: %v\n", err)
						return
					}
					multiBulkLen = c
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

			// execute
			// res, err := parseCommand(elements)
			// if err != nil {
			// 	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
			// } else {
			// 	conn.Write([]byte(res))
			// }

			// could probably reuse existing slices
			buffer = []byte{}
			// elements = []string{}
			// fmt.Println("Parsed: ", elements)
		}
		fmt.Println("Parsed: ", elements)

		out <- elements
	}()

	return out
}
