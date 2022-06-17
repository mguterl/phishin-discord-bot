package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	discordToken := os.Getenv("DISCORD_TOKEN")
	phishinToken := os.Getenv("PHISHIN_TOKEN")
	fmt.Printf("discord: %s\n", discordToken)
	fmt.Printf("phish: %s\n", phishinToken)
}
