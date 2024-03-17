package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
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
	Msgtype byte
	Args    []string
	Error   error
}

// Structured this way to handle packets segmentation of IP protocol
// This thing would be soooo much easier if not segmenting
func Parse(in <-chan []byte) <-chan ParsedMessage {
	out := make(chan ParsedMessage)

	go func() {
		defer close(out)
		var (
			dataType byte
			buf      []byte
			n        int
			args     []string
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
					out <- ParsedMessage{
						Msgtype: SimpleString,
						Args:    []string{},
						Error:   fmt.Errorf("err simple string parsing: %w", err),
					}
					return
				}
				out <- ParsedMessage{
					Msgtype: SimpleString,
					Args:    []string{res},
				}
			case BulkString:
				res, err, _ := parseBulkString(buf, n)
				if err != nil {
					if errors.As(err, &IncompleteMessageError{}) {
						continue
					}
					out <- ParsedMessage{
						Msgtype: SimpleString,
						Args:    []string{},
						Error:   fmt.Errorf("err bulk string parsing: %w", err),
					}
					return
				}
				out <- ParsedMessage{
					Msgtype: Array,
					Args:    []string{res},
				}

			case Array:
				res, err, nn := parseArray(buf, n)
				if err != nil {
					if errors.As(err, &IncompleteMessageError{}) {
						n = nn
						args = append(args, res...)
						continue
					}
					out <- ParsedMessage{
						Msgtype: SimpleString,
						Args:    []string{},
						Error:   fmt.Errorf("err array parsing: %w", err),
					}
					return
				}
				args = append(args, res...)
				out <- ParsedMessage{
					Msgtype: Array,
					Args:    args,
				}
                return
			}
		}

		// dataType = 0
		// buf = []byte{}
		// n = 0
		// args = []string{}
	}()

	return out
}

func parseBulkString(buf []byte, n int) (string, error, int) {
	eol := bytes.IndexByte(buf[n:], '\n')
	if buf[n+eol-1] != '\r' {
		return "", fmt.Errorf("error bulk string protocol parsing: expected \\r\n"), 0
	}
	strLen, err := strconv.Atoi(string(buf[n : n+eol-1]))
	if err != nil {
		return "", fmt.Errorf("error bulk string protocol parsing: invalid length\n"), 0
	}
	if len(buf[n+eol:]) < strLen+2 {
		return "", IncompleteMessageError{msg: "Incomplete message"}, n
	}
	str := buf[n+eol+1 : n+strLen+eol+1]
	if buf[n+eol+strLen+1] != '\r' || buf[n+eol+strLen+2] != '\n' {
		return "", fmt.Errorf("error bulk string protocol parsing: expected \\r\\n\n"), 0
	}
	n += eol + strLen + 3
	return string(str), nil, n
}

func parseArray(buf []byte, n int) ([]string, error, int) {
	var res []string
	eol := bytes.IndexByte(buf[n:], '\n')
	if buf[n+eol-1] != '\r' {
		return res, fmt.Errorf("error array protocol parsing: expected \\r\n"), 0
	}
	numOfEls, err := strconv.Atoi(string(buf[n : eol-1+n]))
	if err != nil {
		return res, fmt.Errorf("error array protocol parsing: invalid length\n"), 0
	}
	n += eol + 1
	for i := 0; i < numOfEls; i++ {
		switch buf[n] {
		case BulkString:
			n++
			str, err, nn := parseBulkString(buf, n)
			n = nn
			if str != "" {
				res = append(res, str)
			}
			if errors.As(err, &IncompleteMessageError{}) {
				return res, err, n
			}
		}
	}

	return res, nil, n
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
