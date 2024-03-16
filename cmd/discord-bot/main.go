package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//go:generate go run github.com/ssouthcity/failsafe/cmd/soyken-tool@latest assets/events

//go:embed assets
var assets embed.FS

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("environment variable DISCORD_TOKEN must be set")
		os.Exit(1)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("invalid client configuration", "err", err)
		os.Exit(1)
	}

	session.AddHandler(setUserStatus)
	session.AddHandler(respondToInteraction)
	session.AddHandler(addImageToEvent)

	if err = session.Open(); err != nil {
		slog.Error("unable to connect to discord ws", "err", err)
		os.Exit(1)
	}
	defer session.Close()

	select {}
}

func setUserStatus(s *discordgo.Session, _ *discordgo.Ready) {
	slog.Info("logged in", "user", s.State.User.Username)

	err := s.UpdateWatchStatus(0, "vex confluxes")
	if err != nil {
		slog.Error("update status failed", "err", err)
	}
}

func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
	if err != nil {
		slog.Error("interaction respond failed", "err", err)
	}
}

func addImageToEvent(session *discordgo.Session, event *discordgo.GuildScheduledEventCreate) {
	tokens := tokenize(event.Name)

	for _, token := range tokens {
		imgPath := fmt.Sprintf("assets/events/%s.soy.jpg", token)

		content, err := assets.ReadFile(imgPath)
		if err != nil {
			continue
		}

		encodedPart := base64.StdEncoding.EncodeToString(content)

		dataURL := "data:image/jpeg;base64," + encodedPart

		_, err = session.GuildScheduledEventEdit(event.GuildID, event.ID, &discordgo.GuildScheduledEventParams{
			Image: dataURL,
		})
		if err != nil {
			slog.Error("scheduled event update failed", "err", err)

			return
		}

		slog.Info("added image to event successfully", "event", event.Name, "image", imgPath)

		return
	}

	slog.Error("did not find any image for event", "event", event.Name)
}

var specialCharRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)

func tokenize(message string) []string {
	lower := strings.ToLower(message)
	alphanums := specialCharRegex.ReplaceAllString(lower, "")
	tokens := strings.Split(alphanums, " ")

	return tokens
}
