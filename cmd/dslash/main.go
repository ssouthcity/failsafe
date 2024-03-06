package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var token = os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("Environment variable DISCORD_TOKEN is unset")
		os.Exit(1)
	}

	var appID = os.Getenv("DISCORD_APPLICATION_ID")
	if appID == "" {
		fmt.Println("Environment variable DISCORD_APPLICATION_ID is unset")
		os.Exit(1)
	}

	var specPath = "application-commands.json"
	if len(os.Args) > 1 {
		specPath = os.Args[1]
	}

	f, err := os.Open(specPath)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer f.Close()

	var commands []*discordgo.ApplicationCommand

	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&commands)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	createdCommands, err := s.ApplicationCommandBulkOverwrite(appID, "", commands)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Printf("The following commands have been created/updated:\n\n")
	for _, cmd := range createdCommands {
		fmt.Printf("  - %s\n", cmd.Name)
	}
}
