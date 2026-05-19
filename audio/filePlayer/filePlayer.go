// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package filePlayer

import (
	"fmt"
	"music-bot/audio/types"
	"os"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "golang.org/x/text/unicode/norm"
	_ "gopkg.in/hraban/opus.v2"
)

type FilePlayer struct {
	isPlaying     bool
	songs         Queue
	timestamp     time.Duration
	currentSong   *Song
	done          chan struct{}
	session       *discordgo.Session
	interaction   *discordgo.InteractionCreate
	connection    *discordgo.VoiceConnection
	isPaused      bool
	needsAdvance  atomic.Bool
}

func (player *FilePlayer) QueueSong(song types.Song) {
	player.songs.Enqueue(song)
	fmt.Println(player.songs.List())
	fmt.Println(player.songs)
}

const delayBetweenSongs = 1 * time.Second

var audioPath = "./audio/filePlayer/files"
var mediaDir string
var cacheDir string

func (player *FilePlayer) TogglePauseResume() {
	//if not on a voice channel, prevent action
	if player.session == nil {
		return
	}
	current := player.CurrentSong()
	if player.IsPlaying() {
		//stop playing
		player.done <- struct{}{}
		player.isPaused = true
		fmt.Println(fmt.Sprintf("Song stopped at %v seconds", player.timestamp.Seconds()))
	} else {
		player.done = make(chan struct{})
		//play the song again (looks at player timestamp)
		fmt.Println(fmt.Sprintf("Resuming song from %v seconds", player.timestamp.Seconds()))
		player.isPaused = false

		//TODO: refactor this to not scatter displaying logic
		//display current song again, when resuming
		player.displayCurrentSong(current)
	}
}

func (player *FilePlayer) CurrentSong() types.Song {
	current := player.songs.Peek()
	return current
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
		mediaDir = audioPath + "/media"
		cacheDir = audioPath + "/cache"
	}
	player.interaction = i
}

func (player *FilePlayer) SetConnection(connection *discordgo.VoiceConnection) {
	player.connection = connection
}

func (player *FilePlayer) IsPlaying() bool {
	return player.isPlaying
}

func (player *FilePlayer) GetQueue() types.Queue {
	return &player.songs
}

func InitFilePlayer() *FilePlayer {
	return &FilePlayer{
		isPlaying:   false,
		songs:       *NewQueue(),
		currentSong: nil,
		done:        make(chan struct{}),
		session:     nil,
		connection:  nil,
		timestamp:   0,
		isPaused:    false,
	}
}

func (player *FilePlayer) Start() {

	for player.songs.Length() > 0 {
		if !player.isPaused {
			current := player.songs.Peek()
			fmt.Println("Current: ", current)
			fmt.Println(player.songs)
			if current != nil {
				player.Play(current.(types.Song))
				time.Sleep(delayBetweenSongs)
				if player.needsAdvance.Load() {
					player.needsAdvance.Store(false)
					player.songs.Dequeue()
				}
			}
		}
	}
}

func (player *FilePlayer) Skip() {
	player.needsAdvance.Store(true)
	player.done <- struct{}{}
	player.timestamp = 0
}

func (player *FilePlayer) Play(song types.Song) {
	go player.displayCurrentSong(song)
	player.done = make(chan struct{})
	vc := player.connection
	player.isPlaying = true

	vc.Speaking(true)
	defer vc.Speaking(false)
	player.streamAudio(vc, song)

}

func (player *FilePlayer) UploadFile(file *discordgo.MessageAttachment) error {
	return saveAttachment(file)
}

// returns the name of the result found
func (player *FilePlayer) FindSong(query *discordgo.ApplicationCommandInteractionDataOption) types.Song {
	return GetClosestMatch(query.StringValue())
}

func (player *FilePlayer) RemoveLastQueued() {
	player.songs.DequeueLastAdded()
}

func (player *FilePlayer) LeaveVoiceChannel() {
	player.done <- struct{}{}
	player.connection.Speaking(false)
	player.session.Close()
}

func (player *FilePlayer) GetUploadedSongs() []string {
	return player.loadFileNames(mediaDir)
}

func (player *FilePlayer) loadFileNames(path string) []string {
	files, _ := os.ReadDir(path)
	fileNames := make([]string, len(files))
	for num, file := range files {
		fileNames[num] = file.Name()
	}
	return fileNames
}

func (player *FilePlayer) SongLength(song types.Song) time.Duration {
	return song.GetDuration()
}

func (player *FilePlayer) SetTimestamp(timestamp time.Duration) {
	player.timestamp = timestamp
}
