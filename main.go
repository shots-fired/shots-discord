package main

import (
	"fmt"
	"os"

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

	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "I am Shots! Hear me roar!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		fmt.Printf("Shots has started on %d servers", len(discord.State.Guilds))
	})

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	defer discord.Close()

	<-make(chan struct{})
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}
