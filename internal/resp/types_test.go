package resp

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_writeBulk(t *testing.T) {
	// given
	input := Value{
		Typ:  BULK.Name,
		Bulk: "Niclas",
	}
	expected := []byte("$6\r\nNiclas\r\n")

	// when
	result := input.Marshal()

	// then
	assert.DeepEqual(t, result, expected)
}

func Test_writeString(t *testing.T) {
	// given
	input := Value{
		Typ: STRING.Name,
		Str: "OK",
	}
	expected := []byte("+OK\r\n")

	// when
	result := input.Marshal()

	// then
	assert.DeepEqual(t, result, expected)
}

func Test_writeError(t *testing.T) {
	// given
	input := Value{
		Typ: ERROR.Name,
		Str: "ERROR",
	}
	expected := []byte("-ERROR\r\n")

	// when
	result := input.Marshal()

	// then
	assert.DeepEqual(t, result, expected)
}

func Test_writeArray(t *testing.T) {
	// given
	input := Value{
		Typ: ARRAY.Name,
		Array: []Value{
			{
				Typ:  BULK.Name,
				Bulk: "Tira",
			},
			{
				Typ:  BULK.Name,
				Bulk: "Misu",
			},
		},
	}
	expected := []byte("*2\r\n$4\r\nTira\r\n$4\r\nMisu\r\n")

	// when
	result := input.Marshal()

	// then
	assert.DeepEqual(t, expected, result)
}

func Test_writeUnknown(t *testing.T) {
	// given
	input := Value{
		Typ: "unknown",
	}

	// when
	result := input.Marshal()

	// then
	assert.Equal(t, len(result), 0)
}
