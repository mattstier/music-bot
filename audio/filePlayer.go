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
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-ar", "48000",
		"-ac", "2",
		"-f", "s16le",
		"pipe:1",
	)
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()

	defer cmd.Wait()

	const frameSize = 960 * 2 * 2 //960 * 2 channels (stereo) * 2 bytes

	// buffered channels
	pcmChannel := make(chan []byte, 50)
	opusChannel := make(chan []byte, 50)

	//buffers
	buf := make([]byte, frameSize)
	data := make([]byte, frameSize)

	//piping pcm
	go func() {
		var pcmBuf []byte
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				pcmBuf = append(pcmBuf, buf[:n]...)
				for len(pcmBuf) >= frameSize {
					frame := make([]byte, frameSize)
					copy(frame, pcmBuf[:frameSize])
					pcmChannel <- frame
					pcmBuf = pcmBuf[frameSize:]
				}
			}
			if err != nil {
				break
			}
		}
		close(pcmChannel)
	}()

	//encoding
	go func() {
		for frame := range pcmChannel {
			samples := bytesToInt16(frame)
			e, err := enc.Encode(samples, data)
			if err != nil {
				break
			}
			pkt := make([]byte, e)
			copy(pkt, data[:e])
			opusChannel <- pkt
		}
		close(opusChannel)
	}()

	time.Sleep(200 * time.Millisecond) // 200 ms latency cushion

	//sending
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for pkt := range opusChannel {
			<-ticker.C
			vc.OpusSend <- pkt
		}
	}()
}
func bytesToInt16(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := range samples {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}
