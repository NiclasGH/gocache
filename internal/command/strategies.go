package command

import (
	"gocache/internal/resp"
	"log"
	"strconv"
	"strings"
)

// / Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.
// / PING {name}?
// / Example:
// / Req: PING
// / Res: PONG
func pingStrategy(args []resp.Value) resp.Value {
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
func setStrategy(args []resp.Value) resp.Value {
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
func getStrategy(args []resp.Value) resp.Value {
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

// / Deletes values at specified keys
// / DEL {key1} [{key2}...]
// / Example:
// / Req: DEL tira
// / Res: (integer) 1
func delStrategy(args []resp.Value) resp.Value {
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

// / Increments number at key. Returns an error if the key is not interpretable as an int
// / INCR {key1}
// / Example:
// / Req: INCR tira
// / Res: (integer) 2
func incrStrategy(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'incr' command"}
	}

	key := args[0].Bulk

	setStorageMutex.Lock()
	defer setStorageMutex.Unlock()

	value, ok := setStorage[key]
	if !ok {
		value = "0"
	}

	savedNumber, err := strconv.Atoi(value)
	if err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: "Value is not a number"}
	}

	savedNumber += 1
	setStorage[key] = strconv.Itoa(savedNumber)

	return resp.Value{Typ: resp.INTEGER.Typ, Num: savedNumber}
}

// / Sets a value in a specific hash at the specified key
// / HSET {hash} {key} {value}
// / Example:
// / Req: HSET tira misu cute
// / Res: OK
func hsetStrategy(args []resp.Value) resp.Value {
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
// / Res: cute
func hgetStrategy(args []resp.Value) resp.Value {
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
func hdelStrategy(args []resp.Value) resp.Value {
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
func hgetAllStrategy(args []resp.Value) resp.Value {
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
func commandMetadataStrategy(args []resp.Value) resp.Value {
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
