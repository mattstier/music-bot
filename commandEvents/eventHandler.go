package commandEvents

import (
	"github.com/bwmarrin/discordgo"
)

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000

func EventListener(s *discordgo.Session, i *discordgo.InteractionCreate) {
	event := i.ApplicationCommandData().Name
	switch event {
	case "play":
		handlePlayEvent(s, i)
	case "skip":
		handleSkipEvent(s, i)
	case "pause":
		handlePauseEvent(s, i)
	}
}

func handlePlayEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	query := i.ApplicationCommandData().Options[0].StringValue()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Playing Song:",
				Description: "\"" + query + "\"",
				//purple color
				Color: PURPLE,
			},
		}},
	})
}

func handleSkipEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Song Skipped",
				Description: "Next song: ...",
				Color:       RED,
			},
		}},
	})
}

func handlePauseEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Song Paused",
				Color: GREEN,
			},
		}},
	})
}
