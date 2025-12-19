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
	"music-bot/audio/filePlayer"
	"music-bot/components"
	"time"

	_ "music-bot/components"

	"github.com/bwmarrin/discordgo"
)

func sendComplex(embeds []*discordgo.MessageEmbed, components []discordgo.MessageComponent, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	})
	fmt.Println(err)
}

func sendEmbed(embeds []*discordgo.MessageEmbed, s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
}

func displaySongQueued(s *discordgo.Session, i *discordgo.InteractionCreate, song string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Song \"" + song + "\" queued",
		Description: i.Member.User.Username + " added a song to the queue",
		Color:       components.PURPLE,
	}
	sendComplex(
		[]*discordgo.MessageEmbed{embed},
		[]discordgo.MessageComponent{components.CancelButton, components.ListQueueButton},
		s, i)
}

func displaySongNotFound(s *discordgo.Session, i *discordgo.InteractionCreate, song string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Song not found",
		Description: fmt.Sprintf("No results for the term %q", song),
		Color:       components.RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}
func displaySongSkipped(s *discordgo.Session, i *discordgo.InteractionCreate, current string) {
	embed := &discordgo.MessageEmbed{
		Title:       i.Member.User.Username + " skipped this song",
		Description: "Playing next song: " + current,
		Color:       components.RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
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
	embed := &discordgo.MessageEmbed{
		Title:       user + " paused this song",
		Description: fmt.Sprintf("\"%v\" (%v)", currentSong, formatTimestamp(position)),
		Color:       components.GREEN,
	}
	sendComplex(
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
	queue := manager.player.GetQueue() //TODO: move this to the parameter once implemented a custom Queue type
	for j := 0; j < len(queue); j++ {
		list += fmt.Sprintf("%d. %s\n", j+1, queue[j])
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Queued songs",
		Description: list,
		Color:       components.PURPLE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayUploadedSongs(s *discordgo.Session, i *discordgo.InteractionCreate, player *filePlayer.FilePlayer) {
	songs := player.GetUploadedSongs()
	var embed discordgo.MessageEmbed
	if len(songs) > 0 {
		list := ""
		for j := 0; j < len(songs); j++ {
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
	sendEmbed([]*discordgo.MessageEmbed{&embed}, s, i)
}
