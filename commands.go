package main

import (
	"github.com/bwmarrin/discordgo"
)

var PlayCommand = &discordgo.ApplicationCommand{
	Name:        "play",
	Description: "Finds and plays a song by name or queues it",
	Options:     playOptions,
}

var playOptions = []*discordgo.ApplicationCommandOption{
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "query",
		Description: "Input the name of the song here",
		Required:    true,
	},
	//TODO: add looping as an optional, now it was buggy for some reason
}

var SkipCommand = &discordgo.ApplicationCommand{
	Name:        "skip",
	Description: "Skips current song, starts playing the next in the queue",
}

var PauseCommand = &discordgo.ApplicationCommand{
	Name:        "pause",
	Description: "Toggle pause/resume current song playing",
}
