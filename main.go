package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var botID string

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_KEY"))
	if err != nil {
		panic(err)
	}

	user, err := discord.User("@me")
	if err != nil {
		panic(err)
	}

	botID = user.ID

	discord.AddHandler(commandHandler)
	discord.AddHandler(readyHandler)

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()

	<-make(chan struct{})
}

func readyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	fmt.Printf("Shots has started on %d servers", len(discord.State.Guilds))
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == botID || message.Author.Bot {
		// Do nothing because the bot is talking
		return
	} else if strings.HasPrefix(message.Content, "!") {
		split := strings.Split(message.Content, " ")
		switch split[0] {
		case "!register":
			registerTwitchHandler(discord, message, split)
		case "!status":
			discord.ChannelMessageSend(message.ChannelID, "Status is WIP")
		}
	}
}

func registerTwitchHandler(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
	if len(splitMessage) > 1 {
		discord.ChannelMessageSend(message.ChannelID, "Registered "+splitMessage[1])
	} else {
		discord.ChannelMessageSend(message.ChannelID, "Need to specify a username to register")
	}
}
