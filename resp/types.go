package resp

type Type byte

// const untyped elements
const (
	Nulls        = '_'
	SimpleString = '+'
	BulkString   = '$'
	Integer      = ':'
	SimpleError  = '-'
	BulkError    = '!'
	Array        = '*'
	Bool         = '#'
	Double       = ','
)

const crlf = "\r\n"

type Envelope struct {
	OpCode  Type
	Integer int
	String  string
	Array   []Envelope
	Double  float64
	Size    int
}
