package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type (
	// Cat is the data model for thecatapi return
	Cat struct {
		URL string `json:"URL"`
	}

	// Cats is a collection of Cat structs
	Cats []Cat
)

func queryCat() (Cat, error) {
	var cats Cats

	baseURL := "https://api.thecatapi.com/v1/images/search?api_key=%s"
	res, err := http.Get(fmt.Sprintf(baseURL, os.Getenv("CATS_API_KEY")))
	if err != nil {
		return cats[0], err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		err := json.NewDecoder(res.Body).Decode(&cats)
		return cats[0], err
	}
	err = errors.New("Recieved: " + res.Status)
	return cats[0], err
}

// catsHandler accepts the cat request from discord
func catsHandler(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
	cat, err := queryCat()
	if err != nil {
		panic(err)
	}
	discord.ChannelMessageSend(message.ChannelID, cat.URL)
}
