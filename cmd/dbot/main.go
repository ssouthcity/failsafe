package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/failsafe/kensoy"
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

	session.AddHandler(kensoy.GuildScheduledEventCreateHandler)

	session.AddHandler(func(s *discordgo.Session, e *discordgo.GuildScheduledEventUpdate) {
		if e.Status != discordgo.GuildScheduledEventStatusActive {
			return
		}

		logger := slog.With("event", e.ID)

		users, err := s.GuildScheduledEventUsers(e.GuildID, e.ID, 100, false, "", "")
		if err != nil {
			logger.Error("unable to get subscribed users of event", "err", err)
		}

		mentions := make([]string, len(users))
		for i, user := range users {
			mentions[i] = user.User.Mention()
		}

		message := fmt.Sprintf("%s is starting!\n%s", e.Name, strings.Join(mentions, " "))

		bountyChannelID := "1216096479416946791"

		_, err = s.ChannelMessageSend(bountyChannelID, message)
		if err != nil {
			logger.Info("unable to get bounty board channel", "err", err)
		}
	})

	// session.AddHandler(func(s *discordgo.Session, e *discordgo.GuildScheduledEventCreate) {
	// 	bountyBoardChannelID := "1216096479416946791"
	// 	eventLink := fmt.Sprintf("https://discord.com/events/%s/%s", e.GuildID, e.ID)

	// 	_, err := s.ChannelMessageSend(bountyBoardChannelID, eventLink)
	// 	if err != nil {
	// 		slog.Error("unable to post event link", "err", err)
	// 	}

	// 	if e.Image != "" {
	// 		return
	// 	}

	// 	key := os.Getenv("BUNGIE_API_KEY")

	// 	resp, err := bungie.SearchEntity(bungie.ActivityDefinition, e.Name, bungie.WithAPIKey(key))
	// 	if err != nil {
	// 		slog.Error("unable to search for activity", "err", err)
	// 		return
	// 	}

	// 	if len(resp.Response.Results.Results) == 0 {
	// 		slog.Info("no activity results found for event", "name", e.Name)
	// 		return
	// 	}

	// 	topResult := resp.Response.Results.Results[0]

	// 	hash2Image := map[uint64]string{
	// 		2906950631: "assets/vow-of-the-disciple.jpg",
	// 		2381413764: "assets/root-of-nightmares.jpg",
	// 	}

	// 	imgpath, ok := hash2Image[topResult.Hash]
	// 	if !ok {
	// 		slog.Info("no images for activity",
	// 			slog.String("name", topResult.DisplayProperties.Name),
	// 			slog.Uint64("hash", topResult.Hash))
	// 		return
	// 	}

	// 	img, err := OpenImage(imgpath)
	// 	if err != nil {
	// 		slog.Error("unable to open activity image", "err", err)
	// 		return
	// 	}

	// 	s.GuildScheduledEventEdit(e.GuildID, e.ID, &discordgo.GuildScheduledEventParams{
	// 		Image: img,
	// 	})
	// })

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
