package command

import (
	"gocache/internal/resp"
	"log"
	"strings"
	"sync"
)

const (
	PING    = "PING"
	SET     = "SET"
	GET     = "GET"
	HSET    = "HSET"
	HGET    = "HGET"
	HGETALL = "HGETALL"
	COMMAND = "COMMAND"
)

type Command = func([]resp.Value) resp.Value

var Commands = map[string]Command{
	PING:    ping,
	SET:     set,
	GET:     get,
	HSET:    hset,
	HGET:    hget,
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
	subCommand := ""
	if len(args) >= 1 {
		subCommand = strings.ToUpper(args[0].Bulk)
	}

	if subCommand == "DOCS" {
		return commandDocs()
	}

	specs := commandSpecs()
	result := make([]resp.Value, 0, len(specs))
	for _, v := range specs {
		if subCommand == "" || v.Command == subCommand {
			result = append(result, v.Value())
		}
	}
	return resp.Value{Typ: resp.ARRAY.Typ, Array: result}
}

func commandSpecs() []resp.CommandSpec {
	return []resp.CommandSpec{
		{
			Command:       "PING",
			ArgCount:      -1,
			Flags:         []string{"readonly", "fast"},
			FirstKey:      1,
			LastKey:       1,
			Steps:         1,
			AclCategories: []string{"@connection", "@fast"},
		},
		{
			Command:       "GET",
			ArgCount:      2,
			Flags:         []string{"readonly", "fast"},
			FirstKey:      1,
			LastKey:       1,
			Steps:         1,
			AclCategories: []string{"@read", "@fast", "@string"},
		},
		{
			Command:       "SET",
			ArgCount:      3,
			Flags:         []string{"write", "fast"},
			FirstKey:      1,
			LastKey:       2,
			Steps:         1,
			AclCategories: []string{"@write", "@slow", "@string"},
		},
		{
			Command:       "HGET",
			ArgCount:      3,
			Flags:         []string{"readonly", "fast"},
			FirstKey:      1,
			LastKey:       2,
			Steps:         1,
			AclCategories: []string{"@read", "@hash", "@fast"},
		},
		{
			Command:       "HSET",
			ArgCount:      4,
			Flags:         []string{"write", "fast"},
			FirstKey:      1,
			LastKey:       3,
			Steps:         1,
			AclCategories: []string{"@write", "@hash", "@fast"},
		},
		{
			Command:       "HGETALL",
			ArgCount:      2,
			Flags:         []string{"readonly"},
			FirstKey:      1,
			LastKey:       1,
			Steps:         1,
			AclCategories: []string{"@read", "@hash", "@slow"},
		},
	}
}

func commandDocs() resp.Value {
	return resp.Value{Typ: "error", Str: "UNIMPLEMENTED"}
}
