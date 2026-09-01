package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//! Load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	//! Read the variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//! Read APP_ENV
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	fmt.Println("Server is Runinng on port: ", port)
	fmt.Printf("Current environtment: %s\n", env)
}
