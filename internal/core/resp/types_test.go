package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeBulk(t *testing.T) {
	// given
	input := Value{
		Typ:  BULK.Typ,
		Bulk: "Niclas",
	}
	expected := []byte("$6\r\nNiclas\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeString(t *testing.T) {
	// given
	input := Value{
		Typ: STRING.Typ,
		Str: "OK",
	}
	expected := []byte("+OK\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeError(t *testing.T) {
	// given
	input := Value{
		Typ: ERROR.Typ,
		Str: "ERROR",
	}
	expected := []byte("-ERROR\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeInteger(t *testing.T) {
	// given
	input := Value{
		Typ: INTEGER.Typ,
		Num: 100,
	}
	expected := []byte(":100\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeArray(t *testing.T) {
	// given
	input := Value{
		Typ: ARRAY.Typ,
		Array: []Value{
			{
				Typ:  BULK.Typ,
				Bulk: "Tira",
			},
			{
				Typ:  BULK.Typ,
				Bulk: "Misu",
			},
		},
	}
	expected := []byte("*2\r\n$4\r\nTira\r\n$4\r\nMisu\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeNull(t *testing.T) {
	// given
	input := Value{
		Typ: "null",
	}
	expected := []byte("$-1\r\n")

	// when
	result := input.Marshal()

	// then
	assert.EqualValues(t, expected, result)
}

func Test_writeUnknown(t *testing.T) {
	// given
	input := Value{
		Typ: "unknown",
	}

	// when
	result := input.Marshal()

	// then
	assert.Equal(t, 0, len(result))
}
