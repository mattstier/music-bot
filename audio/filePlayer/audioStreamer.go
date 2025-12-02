// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

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
)

const frameSize = 960 * channels * 2 //960 * channels * 2 bytes
const sampleRate = 48000
const channels = 2 // 1 for mono; 2 for stereo
const sendRate = 20 * time.Millisecond

func (player *FilePlayer) streamAudio(vc *discordgo.VoiceConnection) {
	//creating pipe with a ffmpeg command
	cmd, cancel := player.startFFMPEG()
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
	//local done channel
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
			//reset timestamp, so next song plays from beginning
			player.timestamp = 0
			fmt.Println("Song finished")
		}
		player.isPlaying = false
		player.done <- struct{}{}

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
			//drop packet if it cannot be sent
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

func (player *FilePlayer) startFFMPEG() (*exec.Cmd, context.CancelFunc) {
	fmt.Println(audioPath)
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-ss", fmt.Sprintf("%.3f", player.timestamp.Seconds()),
		"-i", mediaDir+"/"+player.CurrentSong(),
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-loglevel", "quiet",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", strconv.Itoa(channels),
		"-f", "s16le",
		"pipe:1",
	)
	return cmd, cancel
}
