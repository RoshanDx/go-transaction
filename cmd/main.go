package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-transaction/repository"
	"go-transaction/user"
	"log"
	"os"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load env")
	}

	ctx := context.Background()

	postgresURI := os.Getenv("DATABASE_URL")

	connPool, err := pgxpool.New(ctx, postgresURI)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}
	defer connPool.Close()

	fmt.Println("connected to database")

	store := repository.NewPostgresRepository(connPool)
	userService := user.NewService(store)

	result, err := userService.CreateUser(&user.User{
		Username: "gohan",
		Activate: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created user", result)

	//keep the program running
	//go func() {
	//	for {
	//		fmt.Printf("%v+\n", time.Now())
	//		time.Sleep(10 * time.Second)
	//	}
	//}()
	//
	//quitChannel := make(chan os.Signal, 1)
	//signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	//<-quitChannel
	////time for cleanup before exit
	//fmt.Println("Adios!")
}
