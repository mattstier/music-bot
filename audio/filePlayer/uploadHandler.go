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
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
)

func saveAttachment(file *discordgo.MessageAttachment) error {

	response, err := http.Get(file.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//only save on disk if it is a supported format
	if !formatIsSupported(data) {
		return errors.New("Unsupported file format")
	}
	//make file folder if does not exist already
	os.MkdirAll(mediaDir, 0755)
	return os.WriteFile(mediaDir+"/"+file.Filename, data, 0644)

}

func formatIsSupported(fileContent []byte) bool {
	//make it time out after 2 seconds of not responding
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//ffprobe checks metadata of a file
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error", "-i", "-")
	//stdin pipe is used to send the data to ffmpeg without saving it
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false
	}
	defer stdin.Close()
	go stdin.Write(fileContent)
	err = cmd.Run()
	if err != nil {
		return false
	}
	return true
}
