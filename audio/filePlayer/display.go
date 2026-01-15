// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package filePlayer

import (
	"music-bot/audio/types"
	"music-bot/components"
	_ "music-bot/components"

	"github.com/bwmarrin/discordgo"
)

func (player *FilePlayer) displayCurrentSong(song types.Song) {
	s := player.session
	i := player.interaction
	embed := &discordgo.MessageEmbed{
		Title:       "Playing Song:",
		Description: "\"" + song.GetName() + "\"",
		Color:       components.PURPLE,
	}
	s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					components.PauseButton,
					components.SkipButton,
				},
			},
		},
	})
}
