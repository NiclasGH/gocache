package resp

import (
	"strconv"
)

type Typ struct {
	RespCode byte
	Typ      string
}

var (
	// rd/wrt
	BULK  = Typ{RespCode: '$', Typ: "bulk"}
	ARRAY = Typ{RespCode: '*', Typ: "array"}

	// wrt only
	NULL    = Typ{RespCode: '$', Typ: "null"}
	INTEGER = Typ{RespCode: ':', Typ: "integer"}
	STRING  = Typ{RespCode: '+', Typ: "string"}
	ERROR   = Typ{RespCode: '-', Typ: "error"}
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
	case ARRAY.Typ:
		return v.marshalArray()
	case BULK.Typ:
		return v.marshalBulk()
	case STRING.Typ:
		return v.marshalString()
	case INTEGER.Typ:
		return v.marshalInteger()
	case NULL.Typ:
		return v.marshallNull()
	case ERROR.Typ:
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v Value) GetArgs() []Value {
	return v.Array[1:]
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

func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER.RespCode)
	bytes = append(bytes, []byte(strconv.Itoa(v.Num))...)
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

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
