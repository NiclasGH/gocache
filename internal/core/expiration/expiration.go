package expiration

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"log"
)

/// passive vs active
/// passive: expires a key when its accessed while expired. This is only weakly implemented here.
/// gocache will simply act as if a key is not there if its expired but will wait for the job to actively delete it

// / active: a job that will occassionally check random keys and expire them
func ExpireRandomKeys(amountOfKeys int, db persistence.Database) {
	for range amountOfKeys {
		k, v, ok := db.GetRandomString()

		if !ok {
			break
		}

		log.Println("Checking key: " + k)
		if v.IsExpired() {
			log.Println("Key is expired: " + k)
			db.DeleteAllStrings(delRequest(k), []string{k})
		}
	}
}

func delRequest(key string) resp.Value {
	return resp.Value{
		Typ: resp.ARRAY.Typ,
		Array: []resp.Value{
			{
				Typ:  resp.BULK.Typ,
				Bulk: "DEL",
			},
			{
				Typ:  resp.BULK.Typ,
				Bulk: key,
			},
		},
	}
}
