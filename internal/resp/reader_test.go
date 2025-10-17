package resp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readBulkMissingLineBreak_succeeds(t *testing.T) {
	// given
	input := "$8\r\nTiramisu"
	expected := Value{
		Typ:  "bulk",
		Bulk: "Tiramisu",
	}

	reader := NewReader(strings.NewReader(input))

	// when
	result, err := reader.Read()

	// then
	if err != nil {
		t.Error("Expected non-failure during parsing. Got an error instead: ", err)
	}

	assert.EqualValues(t, &expected, &result)
}

func Test_readBulk(t *testing.T) {
	// given
	input := "$8\r\nTiramisu\r\n"
	expected := Value{
		Typ:  "bulk",
		Bulk: "Tiramisu",
	}

	reader := NewReader(strings.NewReader(input))

	// when
	result, err := reader.Read()

	// then
	if err != nil {
		t.Error("Expected non-failure during parsing. Got an error instead: ", err)
	}

	assert.EqualValues(t, &expected, &result)
}

func Test_readArrayWith2Bulks(t *testing.T) {
	// given
	input := "*2\r\n$4\r\nTira\r\n$4\r\nMisu\r\n"
	expected := Value{
		Typ: "array",
		Array: []Value{
			{
				Typ:  "bulk",
				Bulk: "Tira",
			},
			{
				Typ:  "bulk",
				Bulk: "Misu",
			},
		},
	}

	reader := NewReader(strings.NewReader(input))

	// when
	result, err := reader.Read()

	// then
	if err != nil {
		t.Error("Expected non-failure during parsing. Got an error instead: ", err)
	}

	assert.EqualValues(t, &expected, &result)
}
