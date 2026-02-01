// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package commandEvents

import (
	"context"
	"fmt"
	"music-bot/audio/filePlayer"
	"music-bot/audio/types"
	"time"

	"github.com/bwmarrin/discordgo"
)

var manager *PlayerManager

type PlayerManager struct {
	player          types.Player
	previousMessage *discordgo.Message
}

func ButtonEventListener(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//filters the interactions to button press events only
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	button := i.MessageComponentData().CustomID
	switch button {
	case "button_cancel":
		manager.handleCancelEvent(s, i)
	case "button_skip":
		manager.handleSkipEvent(s, i)
	case "button_pause":
		manager.handlePauseEvent(s, i)
	case "button_jump_10s_forward":
		manager.handleJumpEvent(s, i, "10s")
	case "button_jump_10s_backward":
		manager.handleJumpEvent(s, i, "-10s")
	case "button_list_queue":
		manager.handleListQueueEvent(s, i)
	case "expand_list":
		handleExpandUploadedListEvent(s, i, manager.player.(*filePlayer.FilePlayer))
	case "collapse_list":
		handleListUploadedEvent(s, i, manager.player.(*filePlayer.FilePlayer))
	}
}

func SlashEventListener(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//filters the interaction for slash command events only
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
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
			manager = &PlayerManager{filePlayer.InitFilePlayer(), nil}
		}
	default:
		if manager == nil {
			manager = &PlayerManager{filePlayer.InitFilePlayer(), nil}
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
	case "list":
		manager.handleListQueueEvent(s, i)
	case "uploaded":
		handleListUploadedEvent(s, i, manager.player.(*filePlayer.FilePlayer))
	case "seek":
		manager.handleSeekEvent(s, i, event.Options[0].StringValue())
	case "jump":
		manager.handleJumpEvent(s, i, event.Options[0].StringValue())
	case "quit":
		manager.handleQuitEvent(s, i)
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
	if result != nil {
		queue := manager.player.GetQueue()
		//only display that its queued if it cannot be immediately played
		if queue != nil && queue.Length() > 0 {
			displaySongQueued(s, i, result)
		}
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
	if !manager.player.IsPlaying() && manager.player.GetQueue().Length() > 0 {
		go manager.player.Start()
	}
}

func (manager *PlayerManager) handleSkipEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	go manager.player.Skip()

	//deleting the previous message about it being paused or resumed
	if manager.previousMessage != nil {
		s.ChannelMessageDelete(i.ChannelID, manager.previousMessage.ID)
	}
	displaySongSkipped(s, i, manager.player.CurrentSong())
}

func (manager *PlayerManager) handlePauseEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {

	action := "resumed"
	if manager.player.IsPlaying() {
		action = "paused"
	}
	if action == "paused" {
		displaySongPaused(s, i)
	}

	//deleting the previous message about it being paused or resumed
	if manager.previousMessage != nil {
		s.ChannelMessageDelete(i.ChannelID, manager.previousMessage.ID)
	}
	go manager.player.TogglePauseResume()
	manager.previousMessage = i.Message
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

func (manager *PlayerManager) handleCancelEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if manager.player != nil {
		manager.player.RemoveLastQueued()
		fmt.Println("Cancelling queueing")
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	//deleting message after hitting the cancel button
	if manager.previousMessage != nil {
		s.ChannelMessageDelete(i.ChannelID, manager.previousMessage.ID)
	}
	manager.previousMessage = i.Message
}

func (manager *PlayerManager) handleListQueueEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if manager.player != nil {
		displayQueue(s, i)
		manager.previousMessage = i.Message
	}
}

func handleExpandUploadedListEvent(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer) {
	if manager.player != nil {
		displayUploadedSongs(s, i, player, true)
	}
}

func handleListUploadedEvent(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer) {
	if manager.player != nil {
		//deleting previous message
		displayUploadedSongs(s, i, player, false)
	}
}

func (manager *PlayerManager) handleQuitEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if manager.player != nil {
		displayQuit(s, i)
		manager.player.LeaveVoiceChannel()
	}
}

func (manager *PlayerManager) handleSeekEvent(s *discordgo.Session, i *discordgo.InteractionCreate, userArg string) {
	if manager.player != nil {

		timestamp, err := parseTimeStamp(userArg)
		if err != nil {
			displayInvalidArgument(s, i, userArg, err)
			return
		}
		manager.player.TogglePauseResume()
		manager.player.SetTimestamp(timestamp)
		//cushion to avoid an unsuccessful timestamp setting
		//TODO: make Toggling atomic
		time.Sleep(500 * time.Millisecond)
		manager.player.TogglePauseResume()
		displayJumpedToTimestamp(s, i)
	}
}

func (manager *PlayerManager) handleJumpEvent(s *discordgo.Session, i *discordgo.InteractionCreate, userArg string) {
	if manager.player != nil {
		parsedTimestamp, err := parseTimeStamp(userArg)
		if err != nil {
			displayInvalidArgument(s, i, userArg, err)
			return
		}
		newTimestamp := manager.player.Timestamp() + parsedTimestamp

		manager.player.TogglePauseResume()
		manager.player.SetTimestamp(newTimestamp)
		//cushion to avoid an unsuccessful timestamp setting
		//TODO: make Toggling atomic
		time.Sleep(500 * time.Millisecond)
		manager.player.TogglePauseResume()
		displayJumpedToTimestamp(s, i)

	}
}
