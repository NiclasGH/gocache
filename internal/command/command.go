package command

import (
	"gocache/internal/resp"
	"log"
	"sync"
)

const (
	PING    = "PING"
	SET     = "SET"
	GET     = "GET"
	HSET    = "HSET"
	HGET    = "HGET"
	HGETALL = "HGETALL"
)

var Commands = map[string]func([]resp.Value) resp.Value{
	PING:    ping,
	SET:     set,
	GET:     get,
	HSET:    hset,
	HGET:    hget,
	HGETALL: hgetAll,
}

var setStorage = map[string]string{}
var setStorageMutex = sync.RWMutex{}

var hsetStorage = map[string]map[string]string{}
var hsetStorageMutex = sync.RWMutex{}

var okResponse = resp.Value{Typ: resp.STRING.Typ, Str: "OK"}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: resp.STRING.Typ, Str: args[0].Bulk}
}

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
