package main

import (
	"gocache/internal/core/startup"
	"gocache/internal/infrastructure"
	"gocache/internal/persistence"
	"log"
	"net"
	"os"
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
	defer database.Close()

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
			err := infrastructure.HandleConnection(connection, database)
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

	database := persistence.NewDatabase(nil)

	startup.ReplayCommands(aof, database)

	return database, nil
}
