package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func ready(logger *slog.Logger) func(*discordgo.Session, *discordgo.Ready) {
	logger = logger.With(slog.String("event", "ready"))

	return (func(_ *discordgo.Session, r *discordgo.Ready) {
		logger.Info("Successfully connected to gateway",
			slog.String("username", r.User.Username),
			slog.Int("guilds", len(r.Guilds)),
		)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	token := "Bot " + os.Getenv("DISCORD_TOKEN")

	session, err := discordgo.New(token)
	if err != nil {
		logger.Error("Unable to create Discord client",
			slog.Any("err", err),
		)
		os.Exit(1)
	}

	session.AddHandler(ready(logger))

	if err := session.Open(); err != nil {
		logger.Error("Discord webhook connection failed to connect",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
	defer session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
