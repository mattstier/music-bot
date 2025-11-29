package main

import (
	"fmt"
	"music-bot/audio/filePlayer"
	"music-bot/commandEvents"
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
		panic(err)
	}

	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	//listens to slash commands and decides which one to execute

	session.AddHandler(commandEvents.EventListener)

	err = session.Open()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started")
	//initializes all the commands
	commandEvents.RegisterCommands(session)

	//this is so that the bot is doing non-blocking wait for interrupts, instead of leaving
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	//finally dump the cache (from searching) before shutdown
	filePlayer.DumpCache()
}
