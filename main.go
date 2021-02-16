package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/shots-fired/shots-common/models"
)

var botID string

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

func registerTwitchHandler(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
	if len(splitMessage) > 1 {
		discord.ChannelMessageSend(message.ChannelID, "Registered "+splitMessage[1])
	} else {
		discord.ChannelMessageSend(message.ChannelID, "Need to specify a username to register")
	}
}

func statusTwitchHandler(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
	str := ""
	res, err := http.Get("http://" + os.Getenv("STORE_ADDRESS") + "/streamers")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	streamers := models.Streamers{}
	json.Unmarshal(body, &streamers)
	for _, v := range streamers {
		str += fmt.Sprintf("%s %s %d\n", v.Name, v.Status, v.Viewers)
		log.Printf("%s %s %d\n", v.Name, v.Status, v.Viewers)
	}
	discord.ChannelMessageSend(message.ChannelID, str)
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	lookup := map[string]func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string){}

	lookup["register"] = registerTwitchHandler
	lookup["status"] = statusTwitchHandler
	lookup["bitch"] = func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("no u <@%s>", message.Author.ID))
	}
	lookup["munt"] = func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
		discord.ChannelMessageSend(message.ChannelID, "munt cuffins 4 lyfe")
	}
	var help = "Hi! I'm shots! Here are my commands: \n"
	for k := range lookup {
		help += fmt.Sprintf("!%s\n", k)
	}
	help += "Some commands only available for admins..."
	lookup["shots"] = func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
		discord.ChannelMessageSend(message.ChannelID, help)
	}

	if message.Author.ID == botID || message.Author.Bot {
		// Do nothing because the bot is talking
		return
	} else if strings.HasPrefix(message.Content, "!") {
		splitMessage := strings.Split(message.Content, " ")
		lookup[splitMessage[0][1:]](discord, message, splitMessage[1:])
	}
}
