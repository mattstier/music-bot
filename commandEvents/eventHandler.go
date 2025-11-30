package commandEvents

import (
	"context"
	"fmt"
	"music-bot/audio"
	"music-bot/audio/filePlayer"
	"music-bot/audio/soundCloudPlayer"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
		if manager == nil {
			manager = &PlayerManager{soundCloudPlayer.InitPlayer()}
		}
	case "FileUpload":
		if manager == nil {
			manager = &PlayerManager{filePlayer.InitFilePlayer()}
		}
	default:
		if manager == nil {
			manager = &PlayerManager{filePlayer.InitFilePlayer()}
		}
	}
	manager.player.SetInteraction(i)

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
	fmt.Println(result)
	if result != "" {
		displaySongQueued(s, i, result)
	} else {
		displaySongNotFound(s, i, query.StringValue())
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

	//waiting for voice connection
	vc := <-voice
	//waiting for voice connection to be ready
	for vc.Cond == nil {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Waiting for websocket to open")
	}
	manager.player.SetSession(s)
	manager.player.SetConnection(vc)
	manager.player.QueueSong(result)

	//only autoplay when otherwise not playing and there are songs to play
	if !manager.player.IsPlaying() && len(manager.player.GetQueue()) > 0 {
		go manager.player.Start()
	}
}

func (manager *PlayerManager) handleSkipEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	next := make(chan string)
	go manager.player.Skip(next)
	current := <-next
	displaySongSkipped(s, i, current)
}

func (manager *PlayerManager) handlePauseEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {

	action := "resumed"
	if manager.player.IsPlaying() {
		action = "paused"
	}
	displaySongPaused(s, i, action)
	go manager.player.TogglePauseResume()
}

func handleFileUploadEvent(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer) {
	attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
	attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentID]
	fmt.Println(attachment.Filename)
	err := player.UploadFile(attachment)
	if err == nil {
		displayUpload(s, i, *attachment)
	} else {
		displayUploadError(s, i, *attachment, err)
	}
}
