package handler

import (
	"errors"
	"fmt"
	"gocache/internal/command"
	"gocache/internal/persistence"
	"gocache/internal/resp"
	"io"
	"net"
	"strings"
)

func HandleConnection(connection net.Conn, database persistence.Database) error {
	// the server allows long lived connections with many commands, until the client closes the connection
	for {
		reader := resp.NewReader(connection)
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println(err)
			return err
		}
		fmt.Printf("Received the following Value: %v\n", value)
		writer := resp.NewWriter(connection)

		if err := verifyValueFormat(value); err != nil {
			fmt.Println(err)
			writer.Write(errorValue(err))
			continue
		}

		commandName, err := retrieveCommandName(value)
		if err != nil {
			fmt.Println(err)
			writer.Write(errorValue(err))
			continue
		}

		command, err := retrieveCommand(commandName)
		if err != nil {
			fmt.Println(err)
			writer.Write(errorValue(err))
			continue
		}

		result := command(value.Array[1:])
		if changedDatabase(result, commandName) {
			database.Save(value)
		}

		writer.Write(result)
	}
}

func verifyValueFormat(value resp.Value) error {
	if value.Typ != resp.ARRAY.Typ || len(value.Array) < 1 {
		return errors.New("Command was sent in an invalid format. It needs to be an array")
	}
	return nil
}

func retrieveCommandName(value resp.Value) (string, error) {
	commandValue := value.Array[0]
	if commandValue.Typ != resp.BULK.Typ {
		return "", errors.New("Unable to read command")
	}

	return strings.ToUpper(commandValue.Bulk), nil
}

func retrieveCommand(name string) (func([]resp.Value) resp.Value, error) {
	command, ok := command.Commands[name]
	if !ok {
		return nil, errors.New("Command is unknown")
	}

	return command, nil
}

func changedDatabase(result resp.Value, commandName string) bool {
	return result.Typ != resp.ERROR.Typ && (commandName == command.SET || commandName == command.HSET)
}

func errorValue(err error) resp.Value {
	return resp.Value{Typ: resp.ERROR.Typ, Str: err.Error()}
}
