package soundCloudPlayer

import "github.com/bwmarrin/discordgo"

const ORANGE = 0xff7700

func (player *SoundCloudPlayer) displayCurrentSong() {
	embed := &discordgo.MessageEmbed{
		Title:       "Playing Song:",
		Description: "\"" + player.CurrentSong() + "\"",
		Color:       ORANGE,
	}
	//show song to be played
	player.session.ChannelMessageSendEmbed(player.interaction.ChannelID, embed)
}
