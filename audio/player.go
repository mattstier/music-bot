package audio

import "github.com/bwmarrin/discordgo"

type Player interface {
	//commands
	Start()
	Play(song string, connection *discordgo.VoiceConnection)
	Stop()
	Resume()
	Skip()
	//misc
	CurrentSong() string
	IsPlaying() bool
	SetSession(session *discordgo.Session)
}

//TODO: Add player manager wrapping player function implementations
//for switching between current play strategies (like Soundcloud, file etc), as it would contain a reference to the current player
