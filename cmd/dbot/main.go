package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"

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
		bountyBoardChannelID := "1216096479416946791"
		eventLink := fmt.Sprintf("https://discord.com/events/%s/%s", e.GuildID, e.ID)

		_, err := s.ChannelMessageSend(bountyBoardChannelID, eventLink)
		if err != nil {
			slog.Error("unable to post event link", "err", err)
		}

		if e.Image != "" {
			return
		}

		key := os.Getenv("BUNGIE_API_KEY")

		resp, err := bungie.SearchEntity(bungie.ActivityDefinition, e.Name, bungie.WithAPIKey(key))
		if err != nil {
			slog.Error("unable to search for activity", "err", err)
			return
		}

		if len(resp.Response.Results.Results) == 0 {
			slog.Info("no activity results found for event", "name", e.Name)
			return
		}

		topResult := resp.Response.Results.Results[0]

		hash2Image := map[uint64]string{
			2906950631: "assets/vow-of-the-disciple.jpg",
		}

		imgpath, ok := hash2Image[topResult.Hash]
		if !ok {
			slog.Info("no images for activity",
				slog.String("name", topResult.DisplayProperties.Name),
				slog.Uint64("hash", topResult.Hash))
			return
		}

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
