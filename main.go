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

	splitMessage := strings.Split(message.Content, " ")
	if lookupFunc, ok := lookup[splitMessage[0][1:]]; ok {
		lookupFunc(discord, message, splitMessage[1:])
	}
}

func mentionHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	baseURL := "https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s"
	splitMessage := strings.Split(message.Content, " ")
	words := []string{splitMessage[1] + "es", splitMessage[1] + "s"}
	for _, word := range words {
		res, err := http.Get(fmt.Sprintf(baseURL, word, os.Getenv("DICTIONARY_API_KEY")))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			if strings.Contains(bodyString, "\"fl\":\"verb\"") {
				discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("_%s %s_", word, strings.Join(splitMessage[2:], " ")))
				return
			}
		}
	}
	discord.ChannelMessageSend(message.ChannelID, "wat")
}

func messageHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == botID || message.Author.Bot {
		// Do nothing because the bot is talking
		return
	} else if strings.HasPrefix(message.Content, "!") {
		commandHandler(discord, message)
	} else if strings.HasPrefix(message.Content, fmt.Sprintf("<@!%s> ", botID)) || strings.HasPrefix(message.Content, fmt.Sprintf("<@%s> ", botID)) {
		mentionHandler(discord, message)
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
