package soundCloudPlayer

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SoundCloudPlayer struct {
	isPlaying   bool
	queue       []string // slice of string for now
	timestamp   time.Duration
	currentSong string
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	connection  *discordgo.VoiceConnection
}

var soundCloudClientID string

func InitPlayer() *SoundCloudPlayer {
	soundCloudClientID = os.Getenv("SOUNDCLOUD_CLIENT_ID")
	fmt.Println(soundCloudClientID)
	return &SoundCloudPlayer{
		isPlaying: false,
		queue:     make([]string, 0),
		timestamp: 0,
	}
}

func (player *SoundCloudPlayer) Start() {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) Play(song string) {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) TogglePauseResume() {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) Skip(next chan string) {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) CurrentSong() string {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) IsPlaying() bool {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) Timestamp() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) FindSong(query *discordgo.ApplicationCommandInteractionDataOption) string {
	return getJSONResponse("something")
}

func (player *SoundCloudPlayer) QueueSong(song string) {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) GetQueue() []string {
	//TODO implement me
	panic("implement me")
}

func (player *SoundCloudPlayer) SetSession(session *discordgo.Session) {
	player.session = session
}

func (player *SoundCloudPlayer) SetInteraction(i *discordgo.InteractionCreate) {
	player.interaction = i
}

func (player *SoundCloudPlayer) SetConnection(vc *discordgo.VoiceConnection) {
	player.connection = vc
}
