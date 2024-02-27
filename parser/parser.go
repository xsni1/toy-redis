package parser

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	SimpleString   byte = '+'
	SimpleError    byte = '-'
	Integer        byte = ':'
	BulkString     byte = '$'
	Array          byte = '*'
	Null           byte = '_'
	Boolean        byte = '#'
	Double         byte = ','
	BigNumber      byte = '('
	BulkError      byte = '!'
	VerbatimString byte = '='
	Map            byte = '%'
	Set            byte = '~'
	Push           byte = '>'
)

type IncompleteMessageError struct {
	msg string
}

func (e IncompleteMessageError) Error() string {
	return e.msg
}

type ParsedMessage struct {
	msgtype byte
	args    []string
}

// Structured this way to handle packets segmentation of IP protocol
func Parse(in <-chan []byte) <-chan ParsedMessage {
	out := make(chan ParsedMessage)

	go func() {
		defer close(out)
		var (
			dataType byte
			buf      []byte
			n        int
		)

		for payload := range in {
			buf = append(buf, payload...)
			if dataType == 0 {
				dataType = getMsgType(payload[0])
				n++
			}

			switch dataType {
			case SimpleString:
				res, err, _ := parseSimpleString(buf, n)
				if err != nil {
					if errors.As(err, &IncompleteMessageError{}) {
						continue
					}
					// TODO: Return error to the client
					//       or return it to `out` channel?
					fmt.Printf("err simple string parsing: %v", err)
					return
				}
				out <- ParsedMessage{
					msgtype: SimpleString,
					args:    []string{res},
				}
			}
		}
		// out <- elements
	}()

	return out
}

func parseSimpleString(buf []byte, n int) (string, error, int) {
	// Incomplete packet, wait for next segment
	if len(buf[n:]) >= 2 && (buf[len(buf)-1] != '\n' || buf[len(buf)-2] != '\r') {
		return "", IncompleteMessageError{msg: "Incomplete message"}, 0
	}

	idx := bytes.IndexByte(buf, '\r')
	if idx != len(buf)-2 {
		return "", fmt.Errorf("Parse error - illegal character \\r or \\n"), 0
	}

	return string(buf[n:idx]), nil, idx - n
}

func getMsgType(b byte) byte {
	var dataType byte

	switch b {
	case SimpleString:
		dataType = SimpleString
	case SimpleError:
		dataType = SimpleError
	case Integer:
		dataType = Integer
	case BulkString:
		dataType = BulkString
	case Array:
		dataType = Array
	case Null:
		dataType = Null
	case Boolean:
		dataType = Boolean
	case Double:
		dataType = Double
	case BigNumber:
		dataType = BigNumber
	case BulkError:
		dataType = BulkString
	case VerbatimString:
		dataType = VerbatimString
	case Map:
		dataType = Map
	case Set:
		dataType = Set
	case Push:
		dataType = Push
	}

	return dataType
}
