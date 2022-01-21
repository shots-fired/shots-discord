package commands

import (
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func hasMovePermission(discord *discordgo.Session, authorID string, channelID string) bool {
	permissions, err := discord.UserChannelPermissions(authorID, channelID)
	if err != nil {
		panic(err)
	}

	// Check if the caller at least has permission to move users to the current channel.
	if (permissions & discordgo.PermissionVoiceMoveMembers) != 0 {
		return true
	} else {
		return false
	}
}

func scatterHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.Message

	callerVoiceState, err := session.State.VoiceState(message.GuildID, message.Author.ID)
	if err != nil {
		log.Printf("%s thought they were real slick calling this from a non-voice channel", message.Author.Username)
		return
	}

	// Check if the caller at least has permission to move users to the current channel.
	if !hasMovePermission(session, message.Author.ID, callerVoiceState.ChannelID) {
		log.Printf("%s can't move people from %s", message.Author.Username, callerVoiceState.ChannelID)
		return
	}

	// Find all the other channels they can move users too.
	channels, err := session.GuildChannels(message.GuildID)
	if err != nil {
		panic(err)
	}

	validChannels := []string{}
	for _, channel := range channels {
		if hasMovePermission(session, message.Author.ID, channel.ID) && channel.Bitrate != 0 {
			validChannels = append(validChannels, channel.ID)
		}
	}

	// get all users in same voice chat and asssign them a fun new channel!
	guild, err := session.State.Guild(message.GuildID)
	if err != nil {
		panic(err)
	}
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	for _, userVoiceState := range guild.VoiceStates {
		if userVoiceState.ChannelID != callerVoiceState.ChannelID {
			continue
		}
		channel := validChannels[random.Intn(len(validChannels))]
		err := session.GuildMemberMove(message.GuildID, userVoiceState.UserID, &channel)
		var logmsg string
		if err != nil {
			logmsg = "Couldn't move %s to %s... oh well..."
		} else {
			logmsg = "Moving %s to %s"
		}
		log.Printf(logmsg, userVoiceState.UserID, channel)
	}
}

func momsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.Message

	callerVoiceState, err := session.State.VoiceState(message.GuildID, message.Author.ID)
	if err != nil {
		log.Printf("%s thought they were real slick calling this from a non-voice channel", message.Author.Username)
		return
	}

	guild, err := session.State.Guild(i.GuildID)
	if err != nil {
		panic(err)
	}

	if guild.AfkChannelID == "" {
		log.Printf("No AFK channel for %s", message.GuildID)
		return
	}

	// Check if the caller at least has permission to move members to AFK.
	if !hasMovePermission(session, message.Author.ID, guild.AfkChannelID) {
		log.Printf("%s can't move people to %s", message.Author.Username, callerVoiceState.ChannelID)
		return
	}

	for _, userVoiceState := range guild.VoiceStates {
		if userVoiceState.ChannelID != callerVoiceState.ChannelID {
			continue
		}
		err := session.GuildMemberMove(message.GuildID, userVoiceState.UserID, &guild.AfkChannelID)
		var logmsg string
		if err != nil {
			logmsg = "Couldn't move %s to %s... oh well..."
		} else {
			logmsg = "Moving %s to %s"
		}
		log.Printf(logmsg, userVoiceState.UserID, guild.AfkChannelID)
	}
}
