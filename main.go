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
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
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
}
