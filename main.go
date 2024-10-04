package main

import (
	"log"
)


func main() {
	log.Println("Starting API Server")
	db, err := NewPostgresDatabase()
	
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}


	var server *APIServer =  NewAPIServer(db, ":3000", "v1", "/api")
	server.Start()
}