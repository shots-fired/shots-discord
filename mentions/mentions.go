package mentions

import (
	"fmt"
	"github.com/shots-fired/shots-common/tools"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MentionHandler is the primary dispatcher for @shots mentions.
func MentionHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	splitMessage := strings.Split(message.Content, " ")
	verbed := tools.Verber(strings.Join(splitMessage[1:], " "))
	if verbed != "" {
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("_%s_", verbed))
	} else {
		discord.ChannelMessageSend(message.ChannelID, "wat")
	}
}
