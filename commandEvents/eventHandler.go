package commandEvents

import (
	"context"
	"fmt"
	"music-bot/audio"
	"music-bot/audio/filePlayer"
	"time"

	"github.com/bwmarrin/discordgo"
)

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000

var manager *PlayerManager

type PlayerManager struct {
	player audio.Player
}

func EventListener(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var platform any
	event := i.ApplicationCommandData()

	//find platform TODO: make into a function
	for i := 0; i < len(event.Options); i++ {
		current := event.Options[i].Value
		//TODO: find a better way to do this
		if current == "Youtube" || current == "SoundCloud" || current == "FileUpload" {
			platform = current
		}
	}
	//select the right audio player for the platform
	switch platform {
	case "Youtube":
		fmt.Println("Not yet implemented...")
		fallthrough
	case "SoundCloud":
		fmt.Println("Not yet implemented...")
		fallthrough
	case "FileUpload":
		if manager == nil {
			manager = &PlayerManager{filePlayer.InitFilePlayer()}
		}
	default:
		if manager == nil {
			manager = &PlayerManager{filePlayer.InitFilePlayer()}
		}
	}

	//select which event to handle
	switch event.Name {
	case "play":
		manager.handlePlayEvent(s, i)
	case "skip":
		manager.handleSkipEvent(s, i)
	case "pause":
		manager.handlePauseEvent(s, i)
	case "upload":

		handleFileUploadEvent(s, i, manager.player.(*filePlayer.FilePlayer))
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

func (manager *PlayerManager) handlePlayEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userChannelID, err := findChannel(s, i)
	if err != nil {
		fmt.Println(userChannelID, " | Error finding channel")
		return
	}
	query := i.ApplicationCommandData().Options[0]
	result := manager.player.FindSong(query)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Song \"" + result + "\" queued",
				Description: i.Member.User.Username + " added a song to the queue",
				Color:       PURPLE,
			},
		}},
	})

	//joining with context, deprecated in the new version
	//but the fork for the fix of the audio channel issue is in v26 not v29
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	voice := make(chan *discordgo.VoiceConnection)

	go func() {
		vc, _ := s.ChannelVoiceJoin(ctx, i.GuildID, userChannelID, false, false)
		voice <- vc
	}()

	//waiting for voice connection
	vc := <-voice
	//waiting for voice connection to be ready
	for vc.Cond == nil {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Waiting for websocket to open")
	}
	manager.player.SetSession(s)
	manager.player.SetConnection(vc)
	manager.player.SetInteraction(i)

	if !manager.player.IsPlaying() {
		go manager.player.Start()
	} else {
		//queue result
		manager.player.QueueSong(result)
	}
}

func (manager *PlayerManager) handleSkipEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	next := make(chan string)
	go manager.player.Skip(next)
	current := <-next
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       i.Member.User.Username + " skipped this song",
				Description: "Playing next song: " + current,
				Color:       RED,
			},
		}},
	})
}

func (manager *PlayerManager) handlePauseEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User.Username
	position := manager.player.Timestamp()
	currentSong := manager.player.CurrentSong()
	action := "resumed"
	if manager.player.IsPlaying() {
		action = "paused"
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
			{
				Title:       user + " " + action + " this song",
				Description: fmt.Sprintf("\"%v\" (%v)", currentSong, formatTimestamp(position)),
				Color:       GREEN,
			},
		}},
	})
	go manager.player.TogglePauseResume()
}

func formatTimestamp(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func handleFileUploadEvent(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentID]
	fmt.Println(attachment.Filename)
	err := player.UploadFile(attachment)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Failed to upload file",
					Description: fmt.Sprintf("File: %v \nError: %v", attachment.Filename, err),
					Color:       RED,
				},
			}},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
				{
					Title:       i.Member.User.Username + " uploaded the following song:",
					Description: "\"" + attachment.Filename + "\"",
					Color:       PURPLE,
				},
			}},
		})
	}
}
