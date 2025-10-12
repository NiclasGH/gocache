package command

import (
	"gocache/internal/resp"
	"sync"
)

var Commands = map[string]func([]resp.Value) resp.Value{
	"PING": Ping,
	"SET":  Ping,
	"GET":  Ping,
}

var storage = map[string]string{}
var storageMutex = sync.RWMutex{}

func Ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: resp.STRING.Typ, Str: args[0].Bulk}
}

func Set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	storageMutex.Lock()
	storage[key] = value
	storageMutex.Unlock()

	return resp.Value{Typ: resp.STRING.Typ, Str: "OK"}
}

func Get(args []resp.Value) resp.Value {
	if len(args) == 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	storageMutex.RLock()
	value, ok := storage[key]
	storageMutex.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}
