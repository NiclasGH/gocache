package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"time"
)

func defaultDb() persistence.Database {
	return persistence.NewDatabase(nil)
}

func request(command string, args []resp.Value) resp.Value {
	request := resp.Value{
		Typ: resp.ARRAY.Typ,
		Array: []resp.Value{
			{
				Typ:  resp.BULK.Typ,
				Bulk: command,
			},
		},
	}
	request.Array = append(request.Array, args...)

	return request
}

func isCloseToTimestamp(t1, t2 time.Time, deviation time.Duration) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	return diff <= deviation
}
