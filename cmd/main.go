package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load env")
	}

	postgresURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		log.Panic(err)
	}

	fmt.Println("connected to database")

	// keep the program running
	go forever()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	//time for cleanup before exit
	fmt.Println("Adios!")
}

func forever() {
	for {
		fmt.Printf("%v+\n", time.Now())
		time.Sleep(10 * time.Second)
	}
}
