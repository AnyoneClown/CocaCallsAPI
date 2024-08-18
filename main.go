package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/AnyoneClown/CocaCallsAPI/api"
)

func main() {
	listenAddr := flag.String("listenaddr", ":8080", "The server adress")

	server := api.NewServer(*listenAddr)
	fmt.Printf("Server is running on port %s", *listenAddr)

	log.Fatal(server.Start())
}
