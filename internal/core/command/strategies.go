package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"log"
	"strconv"
	"strings"
)

// / Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.
// / PING {name}?
// / Example:
// / Req: PING
// / Res: PONG
func pingStrategy(request resp.Value, _ persistence.Database) resp.Value {
	args := request.GetArgs()

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
func setStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	err := db.SaveSet(request, key, value)
	if err != nil {
		return resp.Value{Typ: "error", Str: err.Error()}
	}

	return okResponse
}

// / Gets a value at a specific key
// / GET {key}
// / Example:
// / Req: GET tira
// / Res: misu
func getStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	value, err := db.GetSet(key)
	if err != nil {
		log.Printf("Did not find any value with key %s\n", key)
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

// / Deletes values at specified keys
// / DEL {key1} [{key2}...]
// / Example:
// / Req: DEL tira
// / Res: (integer) 1
func delStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	keys := []string{}
	for _, key := range args {
		if key.Typ != resp.BULK.Typ {
			continue
		}

		keys = append(keys, key.Bulk)
	}

	amountDeleted := db.DeleteAllSet(request, keys)

	return resp.Value{Typ: resp.INTEGER.Typ, Num: amountDeleted}
}

// / Increments number at key. Returns an error if the key is not interpretable as an int
// / INCR {key1}
// / Example:
// / Req: INCR tira
// / Res: (integer) 2
func incrStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'incr' command"}
	}

	key := args[0].Bulk

	value, err := db.GetSet(key)
	if err != nil {
		value = "0"
	}

	savedNumber, err := strconv.Atoi(value)
	if err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: "Value is not a number"}
	}

	savedNumber += 1
	if err = db.SaveSet(request, key, strconv.Itoa(savedNumber)); err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: err.Error()}
	}

	return resp.Value{Typ: resp.INTEGER.Typ, Num: savedNumber}
}

// / Sets a value in a specific hash at the specified key
// / HSET {hash} {key} {value}
// / Example:
// / Req: HSET tira misu cute
// / Res: OK
func hsetStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	if err := db.SaveHSet(request, hash, key, value); err != nil {
		return resp.Value{Typ: "error", Str: err.Error()}
	}

	return okResponse
}

// / Gets a value in a specific hash at the specified key
// / HGET {hash} {key}
// / Example:
// / Req: HGET tira misu
// / Res: cute
func hgetStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	mapValue, err := db.GetHSet(hash)
	if err != nil {
		log.Printf("Did not find any value with hash %s\n", hash)
		return resp.Value{Typ: "null"}
	}
	value, ok := mapValue[key]
	if !ok {
		log.Printf("Did not find any value with key %s\n", key)
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

// / Deletes the specified fields inside a hash
// / HDEL {hash} {key1} [{key2}...]
// / Example:
// / Req: HDEL tira misu
// / Res: (integer) 1
func hdelStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) < 2 {
		return resp.Value{Typ: resp.ERROR.Typ, Str: "ERR wrong number of arguments for 'hdel' command"}
	}

	hashKey := args[0].Bulk

	keys := []string{}
	for _, key := range args[1:] {
		if key.Typ != resp.BULK.Typ {
			continue
		}

		keys = append(keys, key.Bulk)
	}

	amountDeleted := db.DeleteAllHSet(request, hashKey, keys)

	return resp.Value{Typ: resp.INTEGER.Typ, Num: amountDeleted}
}

// / Gets all values of a specific hash
// / HGETALL {hash}
// / Example:
// / Req: HGETALL tira
// / Res:
// / misu
// / cute
func hgetAllStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	value, err := db.GetHSet(hash)

	if err != nil {
		log.Println(err.Error())
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
func commandMetadataStrategy(request resp.Value, _ persistence.Database) resp.Value {
	args := request.GetArgs()

	commandFilter := ""
	if len(args) >= 1 {
		commandFilter = strings.ToUpper(args[0].Bulk)
	}

	var result []resp.Value

	metadata := commandList()
	if commandFilter == "DOCS" {
		if len(args) >= 2 {
			commandFilter = strings.ToUpper(args[1].Bulk)
		} else {
			commandFilter = ""
		}

		result = filterCommands(metadata, commandFilter, (*commandMetadata).docs)
	} else {
		result = filterCommands(metadata, commandFilter, (*commandMetadata).specs)
	}

	return resp.Value{Typ: resp.ARRAY.Typ, Array: result}
}

func filterCommands(items []commandMetadata, filter string, accessor func(*commandMetadata) []resp.Value) []resp.Value {
	result := make([]resp.Value, 0, len(items))
	for i := range items {
		if filter == "" || items[i].name == filter {
			result = append(result, accessor(&items[i])...)
		}
	}
	return result
}

func commandList() []commandMetadata {
	docs := make([]commandMetadata, 0, len(commandMetadatas))
	for _, v := range commandMetadatas {
		docs = append(docs, v)
	}
	return docs
}
