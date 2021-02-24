package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/shots-fired/shots-common/models"
)

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
