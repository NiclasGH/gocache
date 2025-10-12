package resp

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

type Resp struct {
	reader *bufio.Reader
}

func NewReader(input string) *Resp {
	return &Resp{reader: bufio.NewReader(strings.NewReader(input))}
}

func (r *Resp) Read() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		return Value{}, errors.New("Received unknown type: " + string(typ))
	}
}

// type parsers
func (r *Resp) readBulk() (Value, error) {
	v := Value{Typ:"bulk"}

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	r.readLine() // consume last linebreak
	return v, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{Typ:"array"}

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	array := make([]Value, len)
	for i := range len {
		value, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		array[i] = value
	}


	v.Array = array

	r.readLine() // consume last linebreak
	return v, nil
}

// base functions
func (r *Resp) readLine() (line []byte, lineLength int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		lineLength += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], lineLength, nil
}

func (r *Resp) readInteger() (x int, lineLength int, err error) {
	line, lineLength, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, lineLength, err
	}
	return int(i64), lineLength, nil
}
