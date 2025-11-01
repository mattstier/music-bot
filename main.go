package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	//test
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "play" {
			query := i.ApplicationCommandData().Options[0].StringValue()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Song Playing",
						Description: query,
						//purple color
						Color: 0xA21DB9,
					},
				}},
			})
		}
	})
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started")

	//this is so that the bot is doing non-blocking wait for interrupts, instead of leaving
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	//initializes all the commands
	commands := []*discordgo.ApplicationCommand{PlayCommand}
	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			fmt.Println("Could not initialize command", command.Name)
		}
	}
}
