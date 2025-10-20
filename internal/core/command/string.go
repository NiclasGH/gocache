package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"log"
	"strconv"
	"time"
)

// / Saves a value at a specific key
// / SET {key} {value}
// / Example:
// / Req: SET tira misu
// / Res: OK
func setStrategy(request resp.Value, db persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	extraArgs := args[2:]
	expiration := time.Duration(0)
	// We dont check the last one, because we will still need a parameter
	for i := 0; i < len(extraArgs)-1; i++ {
		v := extraArgs[i]

		if v.Typ != resp.BULK.Typ {
			continue
		}
		var unit time.Duration
		switch v.Bulk {
		case "EX":
			unit = time.Second
		case "PX":
			unit = time.Millisecond
		default:
			continue
		}

		rawExpire := extraArgs[i+1]
		expire, err := strconv.Atoi(rawExpire.Bulk)
		if err != nil {
			return resp.Value{Typ: resp.ERROR.Typ, Str: "EX/PX parameter needs to be a number"}
		}

		expiration = unit * time.Duration(expire)
		break
	}

	// for remainingArgs:
	// if remainingArg is known EX or PX
	// take next arg and use time.Second or time.Millisecond

	err := db.SaveString(request, key, persistence.NewString(value, expiration))
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

	value, err := db.GetString(key)
	if err != nil || value.IsExpired() {
		log.Printf("Did not find any value with key %s\n", key)
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value.Value}
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

	amountDeleted, err := db.DeleteAllStrings(request, keys)
	if err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: err.Error()}
	}

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

	value, err := db.GetString(key)
	if err != nil || value.IsExpired() {
		value = persistence.NewString("0", 0)
	}

	savedNumber, err := strconv.Atoi(value.Value)
	if err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: "Value is not a number"}
	}

	savedNumber += 1
	value.SetValue(strconv.Itoa(savedNumber))
	if err = db.SaveString(request, key, value); err != nil {
		return resp.Value{Typ: resp.ERROR.Typ, Str: err.Error()}
	}

	return resp.Value{Typ: resp.INTEGER.Typ, Num: savedNumber}
}
