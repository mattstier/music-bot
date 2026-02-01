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
	Skip()
	CurrentSong() Song
	SongLength(song Song) time.Duration
	IsPlaying() bool
	Timestamp() time.Duration
	FindSong(query *discordgo.ApplicationCommandInteractionDataOption) Song
	QueueSong(song Song)
	GetQueue() Queue
	RemoveLastQueued()
	SetTimestamp(timestamp time.Duration)

	//misc
	SetSession(session *discordgo.Session)
	SetInteraction(i *discordgo.InteractionCreate)
	SetConnection(vc *discordgo.VoiceConnection)
	LeaveVoiceChannel()
}
