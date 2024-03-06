package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
)

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

	session.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		slog.Info("logged in", "user", s.State.User.Username)

		err := s.UpdateWatchStatus(0, "vex confluxes")
		if err != nil {
			slog.Error("update status failed", "err", err)
		}
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "pong",
			},
		})
	})

	if err = session.Open(); err != nil {
		slog.Error("unable to connect to discord ws", "err", err)
		os.Exit(1)
	}
	defer session.Close()

	select {}
}
