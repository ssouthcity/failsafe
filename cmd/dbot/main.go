package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
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

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()

		switch data.Name {
		case "ping":
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong",
				},
			})
			if err != nil {
				slog.Error("interaction respond failed", "err", err)
			}
		case "bounty":
			voicechan, err := s.State.Channel("1175271445174173849")
			if err != nil {
				slog.Error("unable to get voice channel", "err", err)
				return
			}

			timestart := time.Now().Add(time.Minute)

			img, err := OpenImage("assets/vow-of-the-disciple.jpg")
			if err != nil {
				slog.Error("unable to open image", "err", err)
				return
			}

			_, err = s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
				Name:               "Test",
				EntityType:         discordgo.GuildScheduledEventEntityTypeVoice,
				ChannelID:          voicechan.ID,
				ScheduledStartTime: &timestart,
				PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
				Image:              img,
			})
			if err != nil {
				slog.Error("interaction respond failed", "err", err)
			}
		}
	})

	if err = session.Open(); err != nil {
		slog.Error("unable to connect to discord ws", "err", err)
		os.Exit(1)
	}
	defer session.Close()

	select {}
}
