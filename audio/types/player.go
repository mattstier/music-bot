package audio

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Player interface {
	//commands
	Start()
	Play(song string)
	TogglePauseResume()
	Skip(next chan string)
	CurrentSong() string
	CurrentSongLength() (time.Duration, error)
	IsPlaying() bool
	Timestamp() time.Duration
	FindSong(query *discordgo.ApplicationCommandInteractionDataOption) string
	QueueSong(song string)
	GetQueue() []string
	RemoveLastQueued()

	//misc
	SetSession(session *discordgo.Session)
	SetInteraction(i *discordgo.InteractionCreate)
	SetConnection(vc *discordgo.VoiceConnection)
	LeaveVoiceChannel()
}
