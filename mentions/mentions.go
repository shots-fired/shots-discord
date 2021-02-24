package mentions

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MentionHandler is the primary dispatcher for @shots mentions.
func MentionHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
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
