package startup

import (
	"errors"
	"gocache/internal/core/command"
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"strings"
)

func ReplayCommands(disk persistence.DiskPersistence, db persistence.Database) error {
	commands, err := disk.ReadPersistedCommands()
	if err != nil {
		return err
	}

	for _, v := range commands {
		name := strings.ToUpper(v.Array[0].Bulk)
		strategy, ok := command.Strategies[name]
		if !ok {
			return errors.New("Command not found: " + name)
		}

		result := strategy(v, db)
		if result.Typ == resp.ERROR.Typ {
			return errors.New("Command returned error: " + result.Str)
		}
	}

	db.EnablePersistence(disk)

	return nil
}
