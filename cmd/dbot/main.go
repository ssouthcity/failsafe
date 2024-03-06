package main

import (
	"fmt"
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

	auth_scheme := fmt.Sprintf("Bot %s", token)

	s, err := discordgo.New(auth_scheme)
	if err != nil {
		slog.Error("invalid client configuration", "err", err)
		os.Exit(1)
	}

	s.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		slog.Info("logged in", "user", s.State.User.Username)

		err := s.UpdateWatchStatus(0, "vex confluxes")
		if err != nil {
			slog.Error("update status failed", "err", err)
		}
	})

	if err = s.Open(); err != nil {
		slog.Error("unable to connect to discord ws", "err", err)
		os.Exit(1)
	}
	defer s.Close()

	select {}
}
