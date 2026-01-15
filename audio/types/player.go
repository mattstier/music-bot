package types

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Player interface {
	//commands
	Start()
	Play(song Song)
	TogglePauseResume()
	Skip(next chan Song)
	CurrentSong() Song
	SongLength(song Song) time.Duration
	IsPlaying() bool
	Timestamp() time.Duration
	FindSong(query *discordgo.ApplicationCommandInteractionDataOption) Song
	QueueSong(song Song)
	GetQueue() Queue
	RemoveLastQueued()

	//misc
	SetSession(session *discordgo.Session)
	SetInteraction(i *discordgo.InteractionCreate)
	SetConnection(vc *discordgo.VoiceConnection)
	LeaveVoiceChannel()
}
