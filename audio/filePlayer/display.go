// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

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
