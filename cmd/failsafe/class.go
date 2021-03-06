package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/dgimux"
)

func classCommand(config *Config) dgimux.InteractionHandlerFunc {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rw := dgimux.NewResponseWriter()
		rw.Text("Pick your main class")
		rw.Ephemral()
		rw.ComponentRow(discordgo.SelectMenu{
			CustomID:  "class_select",
			MinValues: NewIntPtr(1),
			MaxValues: 1,
			Options: []discordgo.SelectMenuOption{
				{
					Label:       "Titan",
					Value:       "titan",
					Description: "Crayons have been identified in the specimens fists",
					Emoji:       discordgo.ComponentEmoji{ID: "862064884593328199"},
					Default:     ListContainsStr(i.Member.Roles, config.ClassRoles["titan"]),
				},
				{
					Label:       "Hunter",
					Value:       "hunter",
					Description: "Target appears to be moving spastically",
					Emoji:       discordgo.ComponentEmoji{ID: "862064884619542538"},
					Default:     ListContainsStr(i.Member.Roles, config.ClassRoles["hunter"]),
				},
				{
					Label:       "Warlock",
					Value:       "warlock",
					Description: "Specimen seems to be consuming an explosive",
					Emoji:       discordgo.ComponentEmoji{ID: "862064884702773268"},
					Default:     ListContainsStr(i.Member.Roles, config.ClassRoles["warlock"]),
				},
			},
		})

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}

func classSelect(config *Config) dgimux.InteractionHandlerFunc {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		pick := i.MessageComponentData().Values[0]

		rw := dgimux.NewResponseWriter()
		rw.Type(discordgo.InteractionResponseUpdateMessage)
		rw.Text("[joyful] Welcome aboard guardian [pained] There's plenty of canned food in the kitchen unit")
		rw.Ephemral()
		rw.ClearComponentRows()

		for name, roleID := range config.ClassRoles {
			var err error

			if name == pick {
				err = s.GuildMemberRoleAdd(config.GuildID, i.Member.User.ID, config.ClassRoles[pick])
			} else {
				err = s.GuildMemberRoleRemove(config.GuildID, i.Member.User.ID, roleID)
			}

			if err != nil {
				log.Err(err).Str("user", i.Member.User.Username).Bool("add", name == pick).Msg("role was not managed")
				rw.Type(discordgo.InteractionResponseChannelMessageWithSource)
				rw.Text("[upbeat] Oh no! I was unable to identify the guardian. [depressed] Critical failure as per usual...")
				break
			}
		}

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}
