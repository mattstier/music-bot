package filePlayer

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "golang.org/x/text/unicode/norm"
	_ "gopkg.in/hraban/opus.v2"
)

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

func (player *FilePlayer) QueueSong(song string) {
	player.songs = append(player.songs, song)
	fmt.Println(player.songs)
}

const delayBetweenSongs = 1 * time.Second

var audioPath = "./audio/filePlayer/files"

func (player *FilePlayer) TogglePauseResume() {
	//if not on a voice channel, prevent action
	if player.session == nil {
		return
	}
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
	if player.interaction == nil {
		audioPath = audioPath + "/" + i.GuildID
	}
	player.interaction = i
}

func (player *FilePlayer) SetConnection(connection *discordgo.VoiceConnection) {
	player.connection = connection
}

func (player *FilePlayer) IsPlaying() bool {
	return player.isPlaying
}

func (player *FilePlayer) GetQueue() []string {
	return player.songs
}

func InitFilePlayer() *FilePlayer {
	return &FilePlayer{
		isPlaying:   false,
		songs:       make([]string, 0),
		currentSong: 0,
		done:        make(chan struct{}),
		session:     nil,
		connection:  nil,
		timestamp:   0,
	}
}

func (player *FilePlayer) Start() {

	for i := 0; i < len(player.songs); i++ {
		go player.displayCurrentSong()
		current := player.CurrentSong()
		fmt.Println("Current: ", current)
		fmt.Println(player.songs)
		player.Play(current)
		//wait a second between songs
		time.Sleep(delayBetweenSongs)
		//only increment song if the stop is from a skip
		//(only happens when timestamp is zero)
		if player.timestamp == 0 {
			player.currentSong++
		}
	}
}

func (player *FilePlayer) Skip(next chan string) {
	player.currentSong++
	current := player.CurrentSong()
	next <- current

	//stop channel
	player.done <- struct{}{}
	player.timestamp = 0
	go player.Play(current)
}

func (player *FilePlayer) Play(song string) {
	player.done = make(chan struct{})
	vc := player.connection
	player.isPlaying = true

	vc.Speaking(true)
	defer vc.Speaking(false)
	player.streamAudio(vc)

}

func (player *FilePlayer) UploadFile(file *discordgo.MessageAttachment) error {
	return saveAttachment(file)
}

// returns the name of the result found
func (player *FilePlayer) FindSong(query *discordgo.ApplicationCommandInteractionDataOption) string {
	return GetClosestMatch(query.StringValue())
}
