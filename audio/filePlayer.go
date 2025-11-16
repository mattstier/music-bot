package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
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
const audioPath = "./audio/files"
const PURPLE = 0xA21DB9

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
	embed := &discordgo.MessageEmbed{
		Title:       "Playing Song:",
		Description: "\"" + player.CurrentSong() + "\"",
		Color:       PURPLE,
	}
	//show song to be played
	player.session.ChannelMessageSendEmbed(player.interaction.ChannelID, embed)
	vc := player.connection
	player.isPlaying = true

	err := vc.Speaking(true)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.3f", player.timestamp.Seconds()),
		"-i", audioPath+"/"+song,
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", strconv.Itoa(channels),
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
	//sending/streaming pcm into the pcm channel
	go bufferPCM(stdout, pcmChannel)
	//encoding and sending pcms into opus frames to the opus channel
	go encodePCM(pcmChannel, opusChannel)
	// 200 ms latency cushion
	time.Sleep(200 * time.Millisecond)
	//sending opus frames to the VoiceConnection
	go sendOpus(opusChannel, vc, player) //waiting for done to be true

	<-player.done
	player.isPlaying = false
	fmt.Println("Music ended")

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
func sendOpus(opusChannel chan []byte, vc *discordgo.VoiceConnection, player *FilePlayer) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for pkt := range opusChannel {
		<-ticker.C
		vc.OpusSend <- pkt
		player.timestamp += time.Millisecond * 20
	}
	// if no more packets to send, signal done
	player.done <- struct{}{}
	//reset timestamp, so next song plays from beginning
	player.timestamp = 0
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
