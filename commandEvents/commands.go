package commandEvents

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const GUILD_ID = "771489027740139531"

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
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "platform",
		Description: "What platform you want to get the given song from.",
		Choices:     platformSelectMenu,
		Required:    false,
	},
}

var platformSelectMenu = []*discordgo.ApplicationCommandOptionChoice{
	{Name: "Youtube", Value: "Youtube"},
	{Name: "SoundCloud", Value: "SoundCloud"},
	{Name: "File Upload", Value: "FileUpload"},
}

var SkipCommand = &discordgo.ApplicationCommand{
	Name:        "skip",
	Description: "Skips current song, starts playing the next in the queue",
}

var PauseCommand = &discordgo.ApplicationCommand{
	Name:        "pause",
	Description: "Toggle pause/resume current song playing",
}

var UploadFileCommand = &discordgo.ApplicationCommand{
	Name:        "upload",
	Description: "Upload an audio file that you can play later.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionAttachment,
			Name:        "file",
			Description: "File to upload",
			Required:    true,
		},
	},
}

func RegisterCommands(session *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		PlayCommand,
		PauseCommand,
		SkipCommand,
		UploadFileCommand,
	}
	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, GUILD_ID, command)
		if err != nil {
			fmt.Println("Could not initialize command", command.Name)
		}
	}
}
