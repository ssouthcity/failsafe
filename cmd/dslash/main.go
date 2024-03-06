package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var token = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Environment variable DISCORD_TOKEN is unset")
	}

	var appID = os.Getenv("DISCORD_APPLICATION_ID")
	if appID == "" {
		log.Fatal("Environment variable DISCORD_APPLICATION_ID is unset")
	}

	var specPath = "application-commands.json"
	if len(os.Args) > 1 {
		specPath = os.Args[1]
	}

	f, err := os.Open(specPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var commands []*discordgo.ApplicationCommand

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&commands)
	if err != nil {
		log.Fatal(err)
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	createdCommands, err := s.ApplicationCommandBulkOverwrite(appID, "", commands)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("The following commands have been created/updated:")
	for _, cmd := range createdCommands {
		log.Printf("  - %s\n", cmd.Name)
	}
}
