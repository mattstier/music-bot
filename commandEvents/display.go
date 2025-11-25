package commandEvents

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000

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
		Color:       PURPLE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}
func displaySongNotFound(s *discordgo.Session, i *discordgo.InteractionCreate, song string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Song not found",
		Description: fmt.Sprintf("No results for the term %q", song),
		Color:       RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}
func displaySongSkipped(s *discordgo.Session, i *discordgo.InteractionCreate, current string) {
	embed := &discordgo.MessageEmbed{
		Title:       i.Member.User.Username + " skipped this song",
		Description: "Playing next song: " + current,
		Color:       RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayUpload(s *discordgo.Session, i *discordgo.InteractionCreate, attachment discordgo.MessageAttachment) {
	embed := &discordgo.MessageEmbed{
		Title:       i.Member.User.Username + " uploaded the following song:",
		Description: "\"" + attachment.Filename + "\"",
		Color:       PURPLE,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displayUploadError(s *discordgo.Session, i *discordgo.InteractionCreate, attachment discordgo.MessageAttachment, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "Failed to upload file",
		Description: fmt.Sprintf("File: %v \nError: %v", attachment.Filename, err),
		Color:       RED,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func displaySongPaused(s *discordgo.Session, i *discordgo.InteractionCreate, action string) {
	user := i.Member.User.Username
	position := manager.player.Timestamp()
	currentSong := manager.player.CurrentSong()
	embed := &discordgo.MessageEmbed{
		Title:       user + " " + action + " this song",
		Description: fmt.Sprintf("\"%v\" (%v)", currentSong, formatTimestamp(position)),
		Color:       GREEN,
	}
	sendEmbed([]*discordgo.MessageEmbed{embed}, s, i)
}

func formatTimestamp(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
