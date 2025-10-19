package command

import (
	"gocache/internal/core/resp"
)

const (
	PING    = "PING"
	DEL     = "DEL"
	SET     = "SET"
	GET     = "GET"
	INCR    = "INCR"
	HSET    = "HSET"
	HGET    = "HGET"
	HDEL    = "HDEL"
	HGETALL = "HGETALL"
	COMMAND = "COMMAND"
)

var Strategies = map[string]CommandStrategy{
	PING:    pingStrategy,
	SET:     setStrategy,
	GET:     getStrategy,
	DEL:     delStrategy,
	INCR:    incrStrategy,
	HSET:    hsetStrategy,
	HGET:    hgetStrategy,
	HDEL:    hdelStrategy,
	HGETALL: hgetAllStrategy,
	COMMAND: commandMetadataStrategy,
}

var commandMetadatas = []commandMetadata{
	ping,
	get,
	set,
	del,
	incr,
	hget,
	hset,
	hdel,
	hgetAll,
	command,
}

var okResponse = resp.Value{Typ: resp.STRING.Typ, Str: "OK"}
