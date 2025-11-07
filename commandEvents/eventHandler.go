package commandEvents

import (
	"fmt"
	"music-bot/audio"
	"time"

	"github.com/bwmarrin/discordgo"
)

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000

const CHANNEL_ID = "771489027740139536"

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

func findChannel(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	guildID := i.GuildID
	userID := i.Member.User.ID //userID within guild
	voiceState, err := s.State.VoiceState(guildID, userID)
	if err != nil {
		return "", err
	}
	return voiceState.ChannelID, nil
}

func handlePlayEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userChannelID, err := findChannel(s, i)
	if err != nil {
		fmt.Println(userChannelID, " | Error finding channel")
		return
	}
	go s.ChannelVoiceJoin(i.GuildID, userChannelID, false, false)

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
	player := audio.FilePlayer{}
	player.SetSession(s)
	//manually retrieve the current voice connection, because otherwise it fails
	vc := s.VoiceConnections[i.GuildID]
	//wait for the voice connection
	for !vc.Ready {
		time.Sleep(10 * time.Millisecond)
	}
	player.Play(vc)
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
