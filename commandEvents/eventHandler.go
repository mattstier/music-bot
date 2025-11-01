package commandEvents

import (
	"github.com/bwmarrin/discordgo"
)

func EventListener(s *discordgo.Session, i *discordgo.InteractionCreate) {
	event := i.ApplicationCommandData().Name
	switch event {
	case "play":
		handlePlayEvent(s, i)
	case "skip":
		handleSkipEvent(s, i)
	case "pause":
		handlePauseEvent(s, i)
	case "loop":
		handleLoopEvent(s, i)
	}
}

func handlePlayEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	query := i.ApplicationCommandData().Options[0].StringValue()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Playing Song",
				Description: "\"" + query + "\"",
				//purple color
				Color: 0xA21DB9,
			},
		}},
	})
}

func handleSkipEvent(s *discordgo.Session, i *discordgo.InteractionCreate)  {}
func handlePauseEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {}
func handleLoopEvent(s *discordgo.Session, i *discordgo.InteractionCreate)  {}
