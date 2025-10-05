package resp

import (
	"testing"
	"gotest.tools/v3/assert"
)

func Test_parseParam(t *testing.T) {
	// given
	input := "$6\r\nNiclas\r\n"
	expected := Value{
		Typ: "bulk",
		Bulk: "Niclas",
	}

	reader := NewReader(input)

	// when
	result, err := reader.Read()

	// then
	if err != nil {
		t.Error("Expected non-failure during parsing. Got an error instead: ", err)
	}

	assert.DeepEqual(t, &expected, &result)
}
