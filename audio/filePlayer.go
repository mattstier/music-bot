package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
)

const frameSize = 960 * 2 * 2 //960 * 2 channels (stereo) * 2 bytes
const sampleRate = 48000
const channels = 2 // mono; 2 for stereo
const audioPath = "./audio/files"

type FilePlayer struct {
	isPlaying   bool
	songs       []string
	currentSong int
	session     *discordgo.Session
}

func (player *FilePlayer) SetSession(session *discordgo.Session) {
	player.session = session
}
func (player *FilePlayer) IsPlaying() bool {
	return player.isPlaying
}

func (player *FilePlayer) Start(vc *discordgo.VoiceConnection) {
	player.songs = loadFileNames(audioPath)
	for i := player.currentSong; i < len(player.songs); i++ {
		current := player.songs[i]
		fmt.Println("Current: ", current)
		fmt.Println("Position: ", i)
		fmt.Println(player.songs)
		player.Play(current, vc)

		time.Sleep(time.Second)
		player.currentSong++
	}
}

func (player *FilePlayer) Play(song string, vc *discordgo.VoiceConnection) {

	player.isPlaying = true

	err := vc.Speaking(true)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("ffmpeg",
		"-i", audioPath+"/"+song,
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-ar", "48000",
		"-ac", "2",
		"-f", "s16le",
		"pipe:1",
	)
	//creating pipe with a ffmpeg command
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()

	// buffered channels
	pcmChannel := make(chan []byte, 50)
	opusChannel := make(chan []byte, 50)
	done := make(chan bool)

	//sending/streaming pcm into the pcm channel
	go pipePCM(stdout, pcmChannel)

	//encoding and sending pcms into opus frames to the opus channel
	go encodePCM(pcmChannel, opusChannel)

	// 200 ms latency cushion
	time.Sleep(200 * time.Millisecond)

	//sending opus frames to the VoiceConnection
	go sendOpus(opusChannel, done, vc)

	<-done
	player.isPlaying = false
	fmt.Println("Music ended")
}

func pipePCM(stdout io.ReadCloser, pcmChannel chan []byte) {
	buf := make([]byte, frameSize)
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
}

func encodePCM(pcmChannel chan []byte, opusChannel chan []byte) {
	enc, _ := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
	data := make([]byte, frameSize)
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
}

func sendOpus(opusChannel chan []byte, done chan bool, vc *discordgo.VoiceConnection) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for pkt := range opusChannel {
		<-ticker.C
		vc.OpusSend <- pkt
	}
	done <- true
}

func bytesToInt16(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := range samples {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}

func loadFileNames(path string) []string {
	files, _ := os.ReadDir(path)
	fileNames := make([]string, len(files))
	for num, file := range files {
		fileNames[num] = file.Name()
	}
	return fileNames
}
