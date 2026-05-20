// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package main

import (
	"fmt"
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

	session.AddHandler(commandEvents.SlashEventListener)
	session.AddHandler(commandEvents.ButtonEventListener)
	session.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		commandEvents.RegisterGuildCommands(s, g.Guild.ID)
	})

	err = session.Open()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started")
	//initializes all the commands for each guild
	for _, guild := range session.State.Guilds {
		commandEvents.RegisterGuildCommands(session, guild.ID)
	}

	//this is so that the bot is doing non-blocking wait for interrupts, instead of leaving
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
