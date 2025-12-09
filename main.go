package main

import (
	"auth/db"
	"auth/server"
	"log"
	"net/http"
)

func main() {

	db.InitDB()

	m := server.NewServer()

	log.Println("🔥 Сервер аухтится!")
	log.Fatal(http.ListenAndServe(":8081", m))
}