package resp

import "strconv"

type Typ struct {
	RespCode byte
	Name     string
}

var (
	ERROR   = Typ{RespCode: '-', Name: "error"}
	STRING  = Typ{RespCode: '+', Name: "string"}
	INTEGER = Typ{RespCode: ':', Name: "integer"}
	BULK    = Typ{RespCode: '$', Name: "bulk"}
	ARRAY   = Typ{RespCode: '*', Name: "array"}
)

// This can be improved using union types, which go currently do not support
type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case ARRAY.Name:
		return v.marshalArray()
	case BULK.Name:
		return v.marshalBulk()
	case STRING.Name:
		return v.marshalString()
	case ERROR.Name:
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v Value) marshalArray() []byte {
	length := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY.RespCode)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = append(bytes, '\r', '\n')

	for i := range length {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK.RespCode)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING.RespCode)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR.RespCode)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
