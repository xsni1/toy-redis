package parser

import "fmt"

func NullReply() []byte {
	return []byte("_\r\n")
}

func SimpleStringReply(reply string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", reply))
}
