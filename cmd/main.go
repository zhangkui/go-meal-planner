package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zhangkui/go-meal-planner/internal/httpapi"
)

func main() {
	address := os.Getenv("ADDR")
	if address == "" {
		address = ":8080"
	}
	log.Println("meal planner listening on", address)
	if err := http.ListenAndServe(address, httpapi.NewServer().Handler()); err != nil {
		log.Fatal(err)
	}
}
