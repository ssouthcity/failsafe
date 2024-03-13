package kensoy

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func GuildScheduledEventCreateHandler(
	s *discordgo.Session,
	e *discordgo.GuildScheduledEventCreate,
) {
	activity, found := findActivityMentioned(e.Name)

	if !found {
		return
	}

	slog.Info(activity.name, "image", activity.image)
}
