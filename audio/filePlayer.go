package audio

import (
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
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
	const sampleRate = 48000
	const channels = 2 // mono; 2 for stereo

	//opus encoder
	enc, _ := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)

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

	buf := make([]byte, 960*2*2) // 20ms of stereo 16-bit PCM (960 samples * 2 channel * 2 bytes)
	data := make([]byte, 960*2*2)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		n, err := stdout.Read(buf)
		e, err := enc.Encode(bytesToInt16(buf[:n]), data)
		if err != nil {
			break
		}
		<-ticker.C
		vc.OpusSend <- data[:e]
	}
	defer vc.Speaking(false)
}

func bytesToInt16(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := range samples {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}
