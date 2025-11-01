package commandEvents

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const GUILD_ID = ""

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
	{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "loop",
		Description: "The song will be repeated this many times",
		Required:    false,
	},
}

var SkipCommand = &discordgo.ApplicationCommand{
	Name:        "skip",
	Description: "Skips current song, starts playing the next in the queue",
}

var PauseCommand = &discordgo.ApplicationCommand{
	Name:        "pause",
	Description: "Toggle pause/resume current song playing",
}

func RegisterCommands(session *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		PlayCommand,
		PauseCommand,
		SkipCommand,
	}
	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, GUILD_ID, command)
		if err != nil {
			fmt.Println("Could not initialize command", command.Name)
		}
	}
}
