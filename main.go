package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/AnyoneClown/CocaCallsAPI/api"
	"github.com/AnyoneClown/CocaCallsAPI/storage"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "The server adress")

	storage := storage.NewCockroachDB()

	server := api.NewServer(*listenAddr, *storage)
	fmt.Printf("Server is running on port %s", *listenAddr)

	log.Fatal(server.Start())
}
