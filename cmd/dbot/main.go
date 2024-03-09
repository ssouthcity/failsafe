package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/failsafe/bungie"
)

func OpenImage(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	e := base64.StdEncoding.EncodeToString(bs)

	datauri := fmt.Sprintf("data:image/jpeg;base64,%s", e)

	return datauri, nil
}

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

	session.AddHandler(func(s *discordgo.Session, e *discordgo.GuildScheduledEventCreate) {
		imgpath := "assets/vow-of-the-disciple.jpg"

		img, err := OpenImage(imgpath)
		if err != nil {
			slog.Error("unable to open activity image", "err", err)
			return
		}

		s.GuildScheduledEventEdit(e.GuildID, e.ID, &discordgo.GuildScheduledEventParams{
			Image: img,
		})
	})

	// imgpaths := map[string]string{
	// 	"Vow of the Disciple":         "assets/vow-of-the-disciple.jpg",
	// 	"Vow of the Disciple: Normal": "assets/vow-of-the-disciple.jpg",
	// 	"Vow of the Disciple: Master": "assets/vow-of-the-disciple.jpg",
	// 	"Vow of the Disciple: Legend": "assets/vow-of-the-disciple.jpg",
	// 	"Root of Nightmares: Normal":  "assets/root-of-nightmares.jpg",
	// 	"Root of Nightmares: Master":  "assets/root-of-nightmares.jpg",
	// }

	if err = session.Open(); err != nil {
		slog.Error("unable to connect to discord ws", "err", err)
		os.Exit(1)
	}
	defer session.Close()

	select {}
}
