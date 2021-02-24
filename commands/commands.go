package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler is the primary dispatcher for !-prefixed commands in shots.
func CommandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	lookup := map[string]func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string){}

	// lookup["register"] = registerTwitchHandler
	// lookup["status"] = statusTwitchHandler

	// animals
	lookup["cat"] = catsHandler
	lookup["dog"] = dogsHandler

	// fun dispersal commands
	lookup["scatter"] = scatterHandler
	lookup["moms-home"] = momsHandler

	// one-off testers
	lookup["bitch"] = func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("no u <@%s>", message.Author.ID))
	}
	lookup["munt"] = func(discord *discordgo.Session, message *discordgo.MessageCreate, splitMessage []string) {
		discord.ChannelMessageSend(message.ChannelID, "munt cuffins 4 lyfe")
	}

	// constructing the help message
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
