package filePlayer

import (
	"io"
	"net/http"
	"os"

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
	//make file folder if does not exist already
	os.MkdirAll(audioPath, 0755)
	return os.WriteFile(audioPath+"/"+file.Filename, data, 0644)

}
