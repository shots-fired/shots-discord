package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/shots-fired/shots-discord/commands"
	"github.com/shots-fired/shots-discord/mentions"
)

var botID string

func messageHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == botID || message.Author.Bot {
		// Do nothing because the bot is talking
		return
	} else if strings.HasPrefix(message.Content, "!") {
		commands.CommandHandler(discord, message)
	} else if strings.HasPrefix(message.Content, fmt.Sprintf("<@!%s> ", botID)) || strings.HasPrefix(message.Content, fmt.Sprintf("<@%s> ", botID)) {
		mentions.MentionHandler(discord, message)
	}
}

func readyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	fmt.Printf("Shots has started on %d servers", len(discord.State.Guilds))
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	user, err := discord.User("@me")
	if err != nil {
		panic(err)
	}

	botID = user.ID

	discord.AddHandler(messageHandler)
	discord.AddHandler(readyHandler)

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()

	<-make(chan struct{})
}
