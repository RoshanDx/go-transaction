package main

import (
	"fmt"
	"go-transaction/receiver"
	"log"
	"net"
	"net/http"
	"time"
)

func handlerOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world from server one\n")
}

func handlerTwo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Hello, world from server two\n")
}

func main() {

	mux1 := http.NewServeMux()
	mux1.HandleFunc("/serverone", handlerOne)
	httpServer1 := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      mux1,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/servertwo", handlerTwo)
	httpServer2 := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8081"),
		Handler:      mux2,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server1 := receiver.NewServerOne(httpServer1)
	server2 := receiver.NewServerTwo(httpServer2)

	receiverHandler, err := receiver.NewHandler(server1, server2)
	if err != nil {
		log.Fatal(err)
	}

	receiverHandler.Run()

	// -----------------------------------

	//// load .env file
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("unable to load env")
	//}
	//
	//ctx := context.Background()
	//
	//postgresURI := os.Getenv("DATABASE_URL")
	//
	//connPool, err := pgxpool.New(ctx, postgresURI)
	//if err != nil {
	//	log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	//}
	//defer connPool.Close()
	//
	//fmt.Println("connected to database")
	//
	//store := repository.NewPostgresRepository(connPool)
	//userService := user.NewService(store)
	//
	//result, err := userService.CreateUser(&user.User{
	//	Username: "gohan",
	//	Activate: true,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("created user", result)

	// -----------------------------------

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
