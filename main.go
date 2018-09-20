package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var botID string

const commandPrefix string = "!"

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
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	} else if strings.HasPrefix(message.Content, "1") {
		split := strings.Split(message.Content, " ")
		switch split[0] {
		case "!register":
			if len(split) > 1 {
				discord.ChannelMessageSend(message.ChannelID, "Registered "+split[1])
			} else {
				discord.ChannelMessageSend(message.ChannelID, "Need to specify a username to register")
			}
		case "!status":
			discord.ChannelMessageSend(message.ChannelID, "Status is WIP")
		}
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}
