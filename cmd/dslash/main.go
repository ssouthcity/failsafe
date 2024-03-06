package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Environment variable DISCORD_TOKEN is unset")
	}

	appID := os.Getenv("DISCORD_APPLICATION_ID")
	if appID == "" {
		log.Fatal("Environment variable DISCORD_APPLICATION_ID is unset")
	}

	specPath := "application-commands.json"
	if len(os.Args) > 1 {
		specPath = os.Args[1]
	}

	file, err := os.Open(specPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var commands []*discordgo.ApplicationCommand

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&commands)
	if err != nil {
		log.Fatal(err)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	createdCommands, err := session.ApplicationCommandBulkOverwrite(appID, "", commands)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("The following commands have been created/updated:")

	for _, cmd := range createdCommands {
		log.Printf("  - %s\n", cmd.Name)
	}
}
