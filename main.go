package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kctjohnson/jtbubbles/cmd/cli"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cli.Execute()
}
