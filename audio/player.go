package audio

import "github.com/bwmarrin/discordgo"

type Player interface {
	//commands
	Play(connection discordgo.VoiceConnection) error
	Stop() error
	Resume() error
	Skip() error
	//
	SetSession(session *discordgo.Session)
}

//TODO: Add player manager wrapping player function implementations
//for switching between current play strategies (like Soundcloud, file etc), as it would contain a reference to the current player
