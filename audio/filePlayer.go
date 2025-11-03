package audio

import (
	"github.com/bwmarrin/discordgo"
)

var filePlayerInstance *FilePlayer

type FilePlayer struct {
	filePath string
	session  *discordgo.Session
}

// single thread singleton pattern for avoiding bugs on the same thread
func GetFilePlayer() *FilePlayer {
	if filePlayerInstance == nil {
		filePlayerInstance = &FilePlayer{}
	}
	return filePlayerInstance
}

func (player *FilePlayer) SetSession(session *discordgo.Session) {
	player.session = session
}

func (player *FilePlayer) Play() error {
	return nil
}
