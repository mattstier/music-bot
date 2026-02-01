// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package commandEvents

import (
	"fmt"
	"math"
	"music-bot/audio/filePlayer"
	"music-bot/audio/types"
	"music-bot/components"
	"time"

	_ "music-bot/components"

	"github.com/bwmarrin/discordgo"
)

const DisplayListLimit = 5

func displaySongQueued(s *discordgo.Session, i *discordgo.InteractionCreate, song types.Song) {
	embed := &discordgo.MessageEmbed{
		Title:       "Song \"" + song.GetName() + "\" queued",
		Description: i.Member.User.Username + " added a song to the queue",
		Color:       components.PURPLE,
	}
	sendComplexReply(
		[]*discordgo.MessageEmbed{embed},
		[]discordgo.MessageComponent{components.CancelButton, components.ListQueueButton},
		s, i)
}

func displaySongNotFound(s *discordgo.Session, i *discordgo.InteractionCreate, name string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Song not found",
		Description: fmt.Sprintf("No results for the term %q", name),
		Color:       components.RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displaySongSkipped(s *discordgo.Session, i *discordgo.InteractionCreate, current types.Song) {
	var embed *discordgo.MessageEmbed
	if current != nil {
		embed = &discordgo.MessageEmbed{
			Title:       i.Member.User.Username + " skipped this song",
			Description: "Playing next song: " + current.GetName(),
			Color:       components.RED,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       i.Member.User.Username + " skipped this song",
			Description: "No more songs to play. You can add a song to the queue with /play",
			Color:       components.RED,
		}

		sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
	}
}

func displayUpload(s *discordgo.Session, i *discordgo.InteractionCreate, attachment discordgo.MessageAttachment) {
	embed := &discordgo.MessageEmbed{
		Title:       i.Member.User.Username + " uploaded the following song:",
		Description: "\"" + attachment.Filename + "\"",
		Color:       components.BLUE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayUploadError(s *discordgo.Session, i *discordgo.InteractionCreate, attachment discordgo.MessageAttachment, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "Failed to upload file",
		Description: fmt.Sprintf("File: %v \nError: %v", attachment.Filename, err),
		Color:       components.RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayQuit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Quit voice channel",
		Description: "See you next time!",
		Color:       components.PURPLE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displaySongPaused(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User.Username
	position := manager.player.Timestamp()
	currentSong := manager.player.CurrentSong()
	progressBar := generateLoadingBar(position, currentSong.GetDuration(), len(currentSong.GetName()))

	embed := &discordgo.MessageEmbed{
		Title:       user + " paused this song",
		Description: fmt.Sprintf("%s (%s) \n%s", currentSong.GetName(), formatTimestamp(position), progressBar),
		Color:       components.GREEN,
	}
	updateComplexReply(
		[]*discordgo.MessageEmbed{embed},
		[]discordgo.MessageComponent{components.ResumeButton},
		s, i)
}

func formatTimestamp(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func displayQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	list := ""
	queue := manager.player.GetQueue()
	if queue.Length() > 0 {
		position := 1
		for song := range queue.All() {
			name := song.GetName()
			duration := song.GetDuration()
			//stop if they don't have name, this maybe unnecessary
			if name == "" {
				break
			}
			list += fmt.Sprintf("%d. %s (%s)\n", position, name, formatTimestamp(duration))
			position++
		}
	} else {
		list = "There are currently no songs in the queue."
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Queued songs",
		Description: list,
		Color:       components.PURPLE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayUploadedSongs(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer, showAll bool) {
	songs := player.GetUploadedSongs()
	songsToDisplay := len(songs)
	button := components.CollapseListButton
	var embed discordgo.MessageEmbed
	if songsToDisplay > 0 {
		list := ""
		if songsToDisplay >= DisplayListLimit && !showAll {
			songsToDisplay = DisplayListLimit
			button = components.ExpandListButton
		}
		for j := 0; j < songsToDisplay; j++ {
			list += fmt.Sprintf("%d. %s \n", j+1, songs[j])
		}

		embed = discordgo.MessageEmbed{
			Title:       "Uploaded songs available",
			Description: list,
			Color:       components.PURPLE,
		}
	} else {
		embed = discordgo.MessageEmbed{
			Title:       "There are no songs available",
			Description: "Upload a song by writing '/upload'",
			Color:       components.PURPLE,
		}
	}
	// edit message when expanding, send a new one when collapsing/sending it for the first time
	// this assumes the first message is always the unexpanded one
	if showAll {
		updateComplexReply(
			[]*discordgo.MessageEmbed{&embed},
			[]discordgo.MessageComponent{button},
			s, i)
	} else {
		sendComplexReply(
			[]*discordgo.MessageEmbed{&embed},
			[]discordgo.MessageComponent{button},
			s, i)
	}
}

func generateLoadingBar(timestamp time.Duration, songLength time.Duration, size int) string {
	if songLength == 0 {
		return ""
	}
	loadingBar := ""
	conversionRatio := float64(timestamp.Milliseconds()) / float64(songLength.Milliseconds())
	loaded := int(math.Ceil(float64(size) * conversionRatio))

	for i := 0; i < size; i++ {
		if i <= loaded {
			loadingBar += "▓"
		} else {
			loadingBar += "░"
		}
	}
	return loadingBar
}

func displayJumpedToTimestamp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User.Username
	position := manager.player.Timestamp()
	currentSong := manager.player.CurrentSong()
	progressBar := generateLoadingBar(position, currentSong.GetDuration(), len(currentSong.GetName()))

	embed := &discordgo.MessageEmbed{
		Title:       user + " jumped to the following",
		Description: fmt.Sprintf("%s (%s) \n%s", currentSong.GetName(), formatTimestamp(position), progressBar),
		Color:       components.GREEN,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayInvalidArgument(s *discordgo.Session, i *discordgo.InteractionCreate, userArg string, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "Invalid Argument",
		Description: fmt.Sprintf("Argument '%s' is invalid \n Reason: '%s' ", userArg, err),
		Color:       components.RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}
