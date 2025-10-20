// / This test file tests the command documentation, which causes this test to be very long
package command

import (
	"gocache/internal/core/resp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ping(t *testing.T) {
	// given
	expected := resp.Value{Typ: "string", Str: "PONG"}

	ping, ok := Strategies[PING]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	args := request(PING, []resp.Value{})

	// when
	result := ping(args, defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_pingWithArg(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "Tiramisu",
		},
	}

	ping, ok := Strategies[PING]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	expected := resp.Value{Typ: "string", Str: "Tiramisu"}

	// when
	result := ping(request(PING, args), defaultDb())

	// then
	assert.EqualValues(t, expected, result)
}

func Test_command_returnsSpecs(t *testing.T) {
	// given
	expected := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "PING"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "GET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
						{Typ: resp.BULK.Typ, Bulk: "@string"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "SET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
						{Typ: resp.BULK.Typ, Bulk: "@string"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "DEL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
						{Typ: resp.BULK.Typ, Bulk: "@keyspace"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "INCR"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
						{Typ: resp.BULK.Typ, Bulk: "@string"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HGET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HSET"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 4},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 3},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HDEL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -3},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "write"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@write"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "HGETALL"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@read"},
						{Typ: resp.BULK.Typ, Bulk: "@hash"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "COMMAND"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
					},
				},
			},
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "COMMAND DOCS"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -2},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 2},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@slow"},
					},
				},
			},
		},
	}

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, []resp.Value{}), defaultDb())

	// then
	assert.ElementsMatch(t, expected, result.Array)
}

func Test_command_withFilter_caseInsensitive_returnsSpecOfFilter(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PiNg",
		},
	}

	expected := []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				// 1. command
				{Typ: resp.BULK.Typ, Bulk: "PING"},
				// 2. arg count
				{Typ: resp.INTEGER.Typ, Num: -1},
				// 3. flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "readonly"},
						{Typ: resp.BULK.Typ, Bulk: "fast"},
					},
				},
				// 4. first key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 5. last key
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 6. steps between keys
				{Typ: resp.INTEGER.Typ, Num: 1},
				// 7. ACL flags
				{
					Typ: resp.ARRAY.Typ,
					Array: []resp.Value{
						{Typ: resp.BULK.Typ, Bulk: "@connection"},
						{Typ: resp.BULK.Typ, Bulk: "@fast"},
					},
				},
			},
		},
	}

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.ElementsMatch(t, expected, result.Array)
}

func Test_commandDocs_returnsDocs(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DOCS",
		},
	}

	expected := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PING",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "GET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Get the value of key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "string"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "SET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Set key to hold the string value."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "string"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DEL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Removes the specified keys."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "keyspace"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1) - O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "INCR",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Increments the number stored at key by one."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "string"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HGET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns the value associated with field in the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HSET",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Sets the specified fields to their respective values in the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HDEL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Removes the specified fields from the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "keyspace"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "HGETALL",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns all fields and values of the hash stored at key."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "hash"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "COMMAND",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Return an array with details about every Redis command."},
				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "2.8.13"},
				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},
				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "COMMAND DOCS",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Return documentary information about commands."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "7.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(N)"},
			},
		},
	}

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.Equal(t, expected, result.Array)
}

func Test_commandDocs_withFilter_caseInsensitive_returnsDocsOfFilter(t *testing.T) {
	// given
	args := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "DOCS",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PiNg",
		},
	}

	expected := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "PING",
		},
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: "summary"},
				{Typ: resp.BULK.Typ, Bulk: "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk."},

				{Typ: resp.BULK.Typ, Bulk: "since"},
				{Typ: resp.BULK.Typ, Bulk: "1.0.0"},

				{Typ: resp.BULK.Typ, Bulk: "group"},
				{Typ: resp.BULK.Typ, Bulk: "connection"},

				{Typ: resp.BULK.Typ, Bulk: "complexity"},
				{Typ: resp.BULK.Typ, Bulk: "O(1)"},
			},
		},
	}

	commandSpecs, ok := Strategies["COMMAND"]
	if !ok {
		t.Error("Command does not exist")
		return
	}

	// when
	result := commandSpecs(request(COMMAND, args), defaultDb())

	// then
	assert.Equal(t, expected, result.Array)
}
