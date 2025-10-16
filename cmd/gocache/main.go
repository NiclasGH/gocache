package main

import (
	"fmt"
	"gocache/internal/handler"
	"gocache/internal/persistence"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	database, err := initializeDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		// Waits until a message is received. Then returns connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			defer connection.Close()
			handler.HandleConnection(connection, database)
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

	// aof.Initialize()
	return aof, nil
}
