package audio

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"

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
func (player *FilePlayer) Play(vc *discordgo.VoiceConnection) {
	err := vc.Speaking(true)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("ffmpeg",
		"-i", "./audio/music.mp3",
		"-ar", "48000",
		"-ac", "2",
		"-f", "s16le",
		"pipe:1",
	)
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()

	defer cmd.Wait()

	buf := make([]byte, 960*2*2) // 20ms of stereo 16-bit PCM (960 samples * 2 channels * 2 bytes)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		n, err := io.ReadFull(stdout, buf)
		if err != nil {
			break
		}
		<-ticker.C
		vc.OpusSend <- buf[:n]
	}
	defer vc.Speaking(false)
}
