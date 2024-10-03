package main

import (
	"log"
)


func main() {
	log.Println("Starting API Server")
	var server *APIServer =  NewAPIServer(":3000", "v1", "/api")
	server.Start()
}