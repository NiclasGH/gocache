package command

import (
	"gocache/internal/resp"
	"log"
	"strings"
	"sync"
)

const (
	PING    = "PING"
	DEL     = "DEL"
	SET     = "SET"
	GET     = "GET"
	HSET    = "HSET"
	HGET    = "HGET"
	HDEL    = "HDEL"
	HGETALL = "HGETALL"
	COMMAND = "COMMAND"
)

type Command = func([]resp.Value) resp.Value

var Commands = map[string]Command{
	PING:    ping,
	SET:     set,
	DEL:     del,
	GET:     get,
	HSET:    hset,
	HGET:    hget,
	HDEL:    hdel,
	HGETALL: hgetAll,
	COMMAND: command,
}

var setStorage = map[string]string{}
var setStorageMutex = sync.RWMutex{}

var hsetStorage = map[string]map[string]string{}
var hsetStorageMutex = sync.RWMutex{}

var okResponse = resp.Value{Typ: resp.STRING.Typ, Str: "OK"}

// / Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.
// / PING {name}?
// / Example:
// / Req: PING
// / Res: PONG
func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: resp.STRING.Typ, Str: args[0].Bulk}
}

// / Saves a value at a specific key
// / SET {key} {value}
// / Example:
// / Req: SET tira misu
// / Res: OK
func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	setStorageMutex.Lock()
	setStorage[key] = value
	setStorageMutex.Unlock()

	return okResponse
}

// / Deletes values at specified keys
// / DEL {key1} [{key2}...]
// / Example:
// / Req: DEL tira
// / Res: (integer) 1
func del(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	setStorageMutex.Lock()
	defer setStorageMutex.Unlock()

	amountDeleted := 0
	for _, key := range args {
		if key.Typ != resp.BULK.Typ {
			continue
		}

		if _, ok := setStorage[key.Bulk]; ok {
			delete(setStorage, key.Bulk)
			amountDeleted += 1
		}
	}

	return resp.Value{Typ: resp.INTEGER.Typ, Num: amountDeleted}
}

// / Gets a value at a specific key
// / GET {key}
// / Example:
// / Req: GET tira
// / Res: misu
func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	setStorageMutex.RLock()
	value, ok := setStorage[key]
	setStorageMutex.RUnlock()

	if !ok {
		log.Printf("Did not find any value with key %s\n", key)
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

// / Sets a value in a specific hash at the specified key
// / HSET {hash} {key} {value}
// / Example:
// / Req: HSET tira misu cute
// / Res: OK
func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	hsetStorageMutex.Lock()
	if _, ok := hsetStorage[hash]; !ok {
		hsetStorage[hash] = map[string]string{}
	}
	hsetStorage[hash][key] = value
	hsetStorageMutex.Unlock()

	return okResponse
}

// / Gets a value in a specific hash at the specified key
// / HGET {hash} {key}
// / Example:
// / Req: HGET tira misu
// / Res:cute
func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	hsetStorageMutex.RLock()
	value, ok := hsetStorage[hash][key]
	hsetStorageMutex.RUnlock()

	if !ok {
		log.Printf("Did not find any value with hash %s or key %s\n", hash, key)
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

// / Deletes the specified fields inside a hash
// / HDEL {hash} {key1} [{key2}...]
// / Example:
// / Req: HDEL tira misu
// / Res: (integer) 1
func hdel(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR.Typ, Str: "ERR wrong number of arguments for 'hdel' command"}
	}

	hashKey := args[0].Bulk
	hsetStorageMutex.Lock()
	defer hsetStorageMutex.Unlock()

	hash, ok := hsetStorage[hashKey]
	if !ok {
		return resp.Value{Typ: resp.INTEGER.Typ, Num: 0}
	}

	amountDeleted := 0

	for _, key := range args[1:] {
		if key.Typ != resp.BULK.Typ {
			continue
		}
		if _, ok := hash[key.Bulk]; ok {
			delete(hash, key.Bulk)
			amountDeleted++
		}
	}

	if len(hash) == 0 {
		delete(hsetStorage, hashKey)
	}

	return resp.Value{Typ: resp.INTEGER.Typ, Num: amountDeleted}
}

// / Gets all values of a specific hash
// / HGETALL {hash}
// / Example:
// / Req: HGETALL tira
// / Res:
// / misu
// / cute
func hgetAll(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	hsetStorageMutex.RLock()
	value, ok := hsetStorage[hash]
	hsetStorageMutex.RUnlock()

	if !ok {
		log.Printf("Did not find any value with hash %s\n", hash)
		return resp.Value{Typ: "null"}
	}

	values := []resp.Value{}
	for k, v := range value {
		values = append(values, resp.Value{Typ: resp.BULK.Typ, Bulk: k})
		values = append(values, resp.Value{Typ: resp.BULK.Typ, Bulk: v})
	}

	return resp.Value{Typ: "array", Array: values}
}

// / Gives information about commands and about available commands
// / COMMAND -> All available commands and their specs (command structure, acl categories, tips, key specification and subcommands). For simplicity reason, I will implement only the first seven categories
// / COMMAND {command} -> Same as Command but filtered to the command
// / COMMAND DOCS -> Docs about the commands. may include: summary, since redis version, functional group, complexity, doc_flags, arguments. We only use summary, group and complexity
func command(args []resp.Value) resp.Value {
	commandFilter := ""
	if len(args) >= 1 {
		commandFilter = strings.ToUpper(args[0].Bulk)
	}

	var result []resp.Value

	if commandFilter == "DOCS" {
		docs := commandDocs()

		if len(args) >= 2 {
			commandFilter = strings.ToUpper(args[1].Bulk)
		} else {
			commandFilter = ""
		}

		result = filterAndConvert(docs, commandFilter)
	} else {
		specs := commandSpecs()
		result = filterAndConvert(specs, commandFilter)
	}

	return resp.Value{Typ: resp.ARRAY.Typ, Array: result}
}

func filterAndConvert[T intoValue](items []T, filter string) []resp.Value {
	result := make([]resp.Value, 0, len(items))
	for _, item := range items {
		if filter == "" || item.getCommand() == filter {
			result = append(result, item.values()...)
		}
	}
	return result
}

func commandSpecs() []commandSpec {
	return []commandSpec{
		{
			command:       "PING",
			argCount:      -1,
			flags:         []string{"readonly", "fast"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@connection", "@fast"},
		},
		{
			command:       "GET",
			argCount:      2,
			flags:         []string{"readonly", "fast"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@read", "@fast", "@string"},
		},
		{
			command:       "SET",
			argCount:      3,
			flags:         []string{"write", "fast"},
			firstKey:      1,
			lastKey:       2,
			steps:         1,
			aclCategories: []string{"@write", "@slow", "@string"},
		},
		{
			command:       "DEL",
			argCount:      -2,
			flags:         []string{"write"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@write", "@slow", "@keyspace"},
		},
		{
			command:       "INCR",
			argCount:      2,
			flags:         []string{"write", "fast"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@write", "@fast", "@string"},
		},
		{
			command:       "HGET",
			argCount:      3,
			flags:         []string{"readonly", "fast"},
			firstKey:      1,
			lastKey:       2,
			steps:         1,
			aclCategories: []string{"@read", "@hash", "@fast"},
		},
		{
			command:       "HSET",
			argCount:      4,
			flags:         []string{"write", "fast"},
			firstKey:      1,
			lastKey:       3,
			steps:         1,
			aclCategories: []string{"@write", "@hash", "@fast"},
		},
		{
			command:       "HDEL",
			argCount:      -3,
			flags:         []string{"write"},
			firstKey:      1,
			lastKey:       2,
			steps:         1,
			aclCategories: []string{"@write", "@fast", "@hash"},
		},
		{
			command:       "HGETALL",
			argCount:      2,
			flags:         []string{"readonly"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@read", "@hash", "@slow"},
		},
		{
			command:       "COMMAND",
			argCount:      -1,
			flags:         []string{"readonly"},
			firstKey:      1,
			lastKey:       1,
			steps:         1,
			aclCategories: []string{"@connection", "@slow"},
		},
		{
			command:       "COMMAND DOCS",
			argCount:      -2,
			flags:         []string{"readonly"},
			firstKey:      2,
			lastKey:       2,
			steps:         1,
			aclCategories: []string{"@connection", "@slow"},
		},
	}
}

func commandDocs() []commandDoc {
	return []commandDoc{
		{
			command:    "PING",
			summary:    "Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.",
			since:      "1.0.0",
			group:      "connection",
			complexity: "O(1)",
		},
		{
			command:    "GET",
			summary:    "Get the value of key.",
			since:      "1.0.0",
			group:      "string",
			complexity: "O(1)",
		},
		{
			command:    "SET",
			summary:    "Set key to hold the string value.",
			since:      "1.0.0",
			group:      "string",
			complexity: "O(1)",
		},
		{
			command:    "DEL",
			summary:    "Removes the specified keys.",
			since:      "1.0.0",
			group:      "keyspace",
			complexity: "O(1) - O(N)",
		},
		{
			command:    "INCR",
			summary:    "Increments the number stored at key by one.",
			since:      "1.0.0",
			group:      "string",
			complexity: "O(1)",
		},
		{
			command:    "HGET",
			summary:    "Returns the value associated with field in the hash stored at key.",
			since:      "2.0.0",
			group:      "hash",
			complexity: "O(1)",
		},
		{
			command:    "HSET",
			summary:    "Sets the specified fields to their respective values in the hash stored at key.",
			since:      "2.0.0",
			group:      "hash",
			complexity: "O(1)",
		},
		{
			command:    "HDEL",
			summary:    "Removes the specified fields from the hash stored at key.",
			since:      "2.0.0",
			group:      "keyspace",
			complexity: "O(N)",
		},
		{
			command:    "HGETALL",
			summary:    "Returns all fields and values of the hash stored at key.",
			since:      "2.0.0",
			group:      "hash",
			complexity: "O(N)",
		},
		{
			command:    "COMMAND",
			summary:    "Return an array with details about every Redis command.",
			since:      "2.8.13",
			group:      "connection",
			complexity: "O(N)",
		},
		{
			command:    "COMMAND DOCS",
			summary:    "Return documentary information about commands.",
			since:      "7.0.0",
			group:      "connection",
			complexity: "O(N)",
		},
	}
}
