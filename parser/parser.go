package parser

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
		)

		for payload := range in {
			if dataType == 0 {
				switch payload[0] {
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
			}

		}
		// out <- elements
	}()

	return out
}
