package commandEvents

import (
	"context"
	"fmt"
	"music-bot/audio"
	"time"

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

	//joining with context, deprecated in the new version
	//but the fork for the fix of the audio channel issue is in v26 not v29
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	voice := make(chan *discordgo.VoiceConnection)

	go func() {
		vc, _ := s.ChannelVoiceJoin(ctx, i.GuildID, userChannelID, false, false)
		voice <- vc
	}()

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

	//waiting for voice connection
	vc := <-voice
	//waiting for voice connection to be ready
	for vc.Cond == nil {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Waiting for websocket to open")
	}
	if !player.IsPlaying() {
		player.Start(vc)
	} else {
		//queue result
	}
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
