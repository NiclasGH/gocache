package command

var ping commandMetadata = commandMetadata{
	name: PING,
	spec: commandSpec{
		argCount:      -1,
		flags:         []string{"readonly", "fast"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@connection", "@fast"},
	},

	doc: commandDoc{
		summary:    "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.",
		since:      "1.0.0",
		group:      "connection",
		complexity: "O(1)",
	},
}

var set commandMetadata = commandMetadata{
	name: SET,
	spec: commandSpec{
		argCount:      3,
		flags:         []string{"write", "fast"},
		firstKey:      1,
		lastKey:       2,
		steps:         1,
		aclCategories: []string{"@write", "@slow", "@string"},
	},
	doc: commandDoc{
		summary:    "Set key to hold the string value.",
		since:      "1.0.0",
		group:      "string",
		complexity: "O(1)",
	},
}

var get commandMetadata = commandMetadata{
	name: GET,
	spec: commandSpec{
		argCount:      2,
		flags:         []string{"readonly", "fast"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@read", "@fast", "@string"},
	},
	doc: commandDoc{
		summary:    "Get the value of key.",
		since:      "1.0.0",
		group:      "string",
		complexity: "O(1)",
	},
}

var del commandMetadata = commandMetadata{
	name: DEL,
	spec: commandSpec{
		argCount:      -2,
		flags:         []string{"write"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@write", "@slow", "@keyspace"},
	},
	doc: commandDoc{
		summary:    "Removes the specified keys.",
		since:      "1.0.0",
		group:      "keyspace",
		complexity: "O(1) - O(N)",
	},
}

var incr commandMetadata = commandMetadata{
	name: INCR,
	spec: commandSpec{
		argCount:      2,
		flags:         []string{"write", "fast"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@write", "@fast", "@string"},
	},
	doc: commandDoc{
		summary:    "Increments the number stored at key by one.",
		since:      "1.0.0",
		group:      "string",
		complexity: "O(1)",
	},
}

var hset commandMetadata = commandMetadata{
	name: HSET,
	spec: commandSpec{
		argCount:      4,
		flags:         []string{"write", "fast"},
		firstKey:      1,
		lastKey:       3,
		steps:         1,
		aclCategories: []string{"@write", "@hash", "@fast"},
	},
	doc: commandDoc{
		summary:    "Sets the specified fields to their respective values in the hash stored at key.",
		since:      "2.0.0",
		group:      "hash",
		complexity: "O(1)",
	},
}

var hget commandMetadata = commandMetadata{
	name: HGET,
	spec: commandSpec{
		argCount:      3,
		flags:         []string{"readonly", "fast"},
		firstKey:      1,
		lastKey:       2,
		steps:         1,
		aclCategories: []string{"@read", "@hash", "@fast"},
	},
	doc: commandDoc{
		summary:    "Returns the value associated with field in the hash stored at key.",
		since:      "2.0.0",
		group:      "hash",
		complexity: "O(1)",
	},
}

var hdel commandMetadata = commandMetadata{
	name: HDEL,
	spec: commandSpec{
		argCount:      -3,
		flags:         []string{"write"},
		firstKey:      1,
		lastKey:       2,
		steps:         1,
		aclCategories: []string{"@write", "@fast", "@hash"},
	},
	doc: commandDoc{
		summary:    "Removes the specified fields from the hash stored at key.",
		since:      "2.0.0",
		group:      "keyspace",
		complexity: "O(N)",
	},
}

var hgetAll commandMetadata = commandMetadata{
	name: HGETALL,
	spec: commandSpec{
		argCount:      2,
		flags:         []string{"readonly"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@read", "@hash", "@slow"},
	},
	doc: commandDoc{
		summary:    "Returns all fields and values of the hash stored at key.",
		since:      "2.0.0",
		group:      "hash",
		complexity: "O(N)",
	},
}

var command commandMetadata = commandMetadata{
	name: COMMAND,
	subCommands: []commandMetadata{
		{
			name: COMMAND + " DOCS",
			spec: commandSpec{
				argCount:      -2,
				flags:         []string{"readonly"},
				firstKey:      2,
				lastKey:       2,
				steps:         1,
				aclCategories: []string{"@connection", "@slow"},
			},
			doc: commandDoc{
				summary:    "Return documentary information about commands.",
				since:      "7.0.0",
				group:      "connection",
				complexity: "O(N)",
			},
		},
	},
	spec: commandSpec{
		argCount:      -1,
		flags:         []string{"readonly"},
		firstKey:      1,
		lastKey:       1,
		steps:         1,
		aclCategories: []string{"@connection", "@slow"},
	},
	doc: commandDoc{
		summary:    "Return an array with details about every Redis command.",
		since:      "2.8.13",
		group:      "connection",
		complexity: "O(N)",
	},
}
