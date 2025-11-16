package audio

import "github.com/bwmarrin/discordgo"

type Player interface {
	//commands
	Start()
	Play(song string, connection *discordgo.VoiceConnection)
	ToggleStopResume()
	Skip()
	//misc
	CurrentSong() string
	IsPlaying() bool
	SetSession(session *discordgo.Session)
}
