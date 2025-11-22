package filePlayer

import "github.com/bwmarrin/discordgo"

const PURPLE = 0xA21DB9

func (player *FilePlayer) displayCurrentSong() {
	embed := &discordgo.MessageEmbed{
		Title:       "Playing Song:",
		Description: "\"" + player.CurrentSong() + "\"",
		Color:       PURPLE,
	}
	//show song to be played
	player.session.ChannelMessageSendEmbed(player.interaction.ChannelID, embed)
}
