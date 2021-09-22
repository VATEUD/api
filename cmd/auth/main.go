package main

import (
	"auth/internal/app/auth"
	"log"
)

func main() {
	log.Println("Starting...")
	auth.Start()
}
