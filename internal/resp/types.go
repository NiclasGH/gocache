package resp

const (
	ERROR   = '-'
	STRING  = '+'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// This can be improved using union types
type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

