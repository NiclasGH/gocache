package main

import (
	"gocache/internal/command"
	"gocache/internal/handler"
	"gocache/internal/persistence"
	"gocache/internal/resp"
	"log"
	"net"
	"os"
	"strings"
)

var ready chan struct{}

func main() {
	port, ok := os.LookupEnv("GC_PORT")
	if !ok {
		port = "6379"
	}
	port = ":" + port
	log.Printf("Listening on port %v\n", port)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	database, err := initializeDatabase()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if ready != nil {
		close(ready)
	}

	for {
		// Waits until a message is received. Then returns connection
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			defer connection.Close()
			err := handler.HandleConnection(connection, database)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func initializeDatabase() (persistence.Database, error) {
	databasePath, ok := os.LookupEnv("GC_DATABASE_PATH")
	if !ok {
		databasePath = "database.aof"
	}

	aof, err := persistence.NewAof(databasePath)
	if err != nil {
		return nil, err
	}

	err = aof.Initialize(func(value resp.Value) {
		commandName := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		command, ok := command.Strategies[commandName]
		if !ok {
			return
		}

		command(args)
	})
	if err != nil {
		aof.Close()
		return nil, err
	}

	return aof, nil
}
