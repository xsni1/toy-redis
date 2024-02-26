package parser

import (
	"bytes"
	"fmt"
	"strings"
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

// Structured this way to handle packets segmentation of IP protocol
func Parse(in <-chan []byte) <-chan []string {
	out := make(chan []string)

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
				parseSimpleString(buf, n)
			}
		}
		// out <- elements
	}()

	return out
}

// moge zwracac n
// n byloby dodawane w tej funkcji za kazdy sparsowany token
func parseSimpleString(buf []byte, n int) (int, error) {
	if len(buf) >= 2 && (buf[len(buf)-1] != '\n' || buf[len(buf)-2] != '\r') {
		return 0, fmt.Errorf("Parse error - expected \\r\\n termination")
	}

	idx := bytes.IndexByte(buf, '\r')
	if idx != len(buf)-2 {
		return 0, fmt.Errorf("Parse error - illegal character \\r or \\n")
	}

	return 1, nil
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
