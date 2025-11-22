package filePlayer

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/hraban/opus.v2"
	_ "gopkg.in/hraban/opus.v2"
)

const frameSize = 960 * channels * 2 //960 * channels * 2 bytes
const sampleRate = 48000
const channels = 2 // 1 for mono; 2 for stereo
const sendRate = 20 * time.Millisecond
const audioPath = "./audio/filePlayer/files"

type FilePlayer struct {
	isPlaying   bool
	songs       []string
	timestamp   time.Duration
	currentSong int
	done        chan struct{}
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	connection  *discordgo.VoiceConnection
}

func (player *FilePlayer) TogglePauseResume() {
	if player.IsPlaying() {
		//stop playing
		player.done <- struct{}{}
		fmt.Println(fmt.Sprintf("Song stopped at %v seconds", player.timestamp.Seconds()))
	} else {
		player.done = make(chan struct{})
		//play the song again (looks at player timestamp)
		fmt.Println(fmt.Sprintf("Resuming song from %v seconds", player.timestamp.Seconds()))
		go player.Play(player.CurrentSong())
	}
}

func (player *FilePlayer) CurrentSong() string {

	if player.songs != nil && len(player.songs) > 0 && player.currentSong < len(player.songs) {
		return player.songs[player.currentSong]
	}
	return "Couldn't find next song"
}

func (player *FilePlayer) Timestamp() time.Duration {
	return player.timestamp
}

func (player *FilePlayer) SetSession(session *discordgo.Session) {
	player.session = session
}

func (player *FilePlayer) SetInteraction(i *discordgo.InteractionCreate) {
	player.interaction = i
}

func (player *FilePlayer) SetConnection(connection *discordgo.VoiceConnection) {
	player.connection = connection
}

func (player *FilePlayer) IsPlaying() bool {
	return player.isPlaying
}

func InitFilePlayer() *FilePlayer {
	return &FilePlayer{
		isPlaying:   false,
		songs:       loadFileNames(audioPath),
		currentSong: 0,
		done:        make(chan struct{}),
		session:     nil,
		connection:  nil,
		timestamp:   0,
	}
}

func (player *FilePlayer) Start() {

	for i := player.currentSong; i < len(player.songs); i++ {
		current := player.songs[i]
		fmt.Println("Current: ", current)
		fmt.Println(player.songs)
		go player.Play(current)
		//wait for song to finish
		<-player.done
		//wait a second between songs
		time.Sleep(time.Second)
		player.currentSong++
	}
}

func (player *FilePlayer) Skip(next chan string) {
	player.isPlaying = false
	player.currentSong++
	current := player.CurrentSong()
	next <- current

	//stop channel
	player.done <- struct{}{}
	//clear channel to run again
	player.done = make(chan struct{})
	go player.Play(current)
}

func (player *FilePlayer) Play(song string) {

	vc := player.connection
	player.isPlaying = true
	player.displayCurrentSong()
	vc.Speaking(true)

	//creating pipe with a ffmpeg command
	cmd, cancel := player.startFFMPEG()
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
	//local
	sessionDone := make(chan struct{})

	// buffered channels
	pcmChannel := make(chan []byte, 50)
	opusChannel := make(chan []byte, 50)

	//handles stopping, song finishing etc
	go func() {
		select {
		//triggers if done has been signaled
		case <-player.done:
			stdout.Close() // always close the pipe first, to not break it
			cancel()
			fmt.Println("Player stopped")
		case <-sessionDone:
			fmt.Println("Song finished")
		}
		player.isPlaying = false
		//reset timestamp, so next song plays from beginning
		player.timestamp = 0
		vc.Speaking(false)

	}()
	//sending/streaming pcm into the pcm channel
	go bufferPCM(stdout, pcmChannel)
	//encoding and sending pcms into opus frames to the opus channel
	go encodePCM(pcmChannel, opusChannel)
	// 200 ms latency cushion
	time.Sleep(200 * time.Millisecond)
	//sending opus frames to the VoiceConnection
	go sendOpus(opusChannel, vc, player, sessionDone) //waiting for done to be true

}

// returns the name of the result found
func (player *FilePlayer) FindSong(query *discordgo.ApplicationCommandInteractionDataOption) string {
	//returning the original string => only for now
	return query.StringValue()
}

func bufferPCM(stdout io.ReadCloser, pcmChannel chan []byte) {
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

// encodes pcm to Opus format
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

// sends opus frames to discord voice channel
func sendOpus(opusChannel chan []byte, vc *discordgo.VoiceConnection, player *FilePlayer, sessionDone chan struct{}) {
	ticker := time.NewTicker(sendRate)
	defer ticker.Stop()
	for pkt := range opusChannel {
		<-ticker.C
		select {
		case vc.OpusSend <- pkt:
			//send, and increase timestamp if successful
			player.timestamp += sendRate
		case <-player.done:
			return
		default:
			//drop packet if cannot send
		}

	}
	// if no more packets to send, signal done
	close(sessionDone)
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

func (player *FilePlayer) startFFMPEG() (*exec.Cmd, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-ss", fmt.Sprintf("%.3f", player.timestamp.Seconds()),
		"-i", audioPath+"/"+player.CurrentSong(),
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", strconv.Itoa(channels),
		"-f", "s16le",
		"pipe:1",
	)
	return cmd, cancel
}
