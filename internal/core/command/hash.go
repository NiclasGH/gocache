package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"log"
)

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

	if err := db.SaveHash(request, hash, key, value); err != nil {
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

	mapValue, err := db.GetHash(hash)
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

	amountDeleted, err := db.DeleteAllHashKeys(request, hashKey, keys)
	if err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: err.Error()}
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
func hgetAllStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	value, err := db.GetHash(hash)

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
