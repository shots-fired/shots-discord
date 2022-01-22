package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"
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
	member := i.Member
	if member == nil {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Shhhh I'm sleeping.",
			},
		})
		return
	}
	userID := member.User.ID
	guildID := i.GuildID

	callerVoiceState, err := session.State.VoiceState(guildID, userID)
	if err != nil {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> thought they were real slick calling this from a non-voice channel", userID),
			},
		})
		return
	}

	// Check if the caller at least has permission to move users to the current channel.
	if !hasMovePermission(session, userID, callerVoiceState.ChannelID) {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> can't move people from <#%s>", userID, callerVoiceState.ChannelID),
			},
		})
		return
	}

	// get all users in same voice chat and asssign them a fun new channel!
	guild, err := session.State.Guild(guildID)
	if err != nil {
		panic(err)
	}

	validChannels := []string{}
	for _, channel := range guild.Channels {
		if hasMovePermission(session, userID, channel.ID) && channel.Bitrate != 0 {
			validChannels = append(validChannels, channel.ID)
		}
	}

	log.Printf("Can move to %v", validChannels)

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	for _, userVoiceState := range guild.VoiceStates {
		if userVoiceState.ChannelID != callerVoiceState.ChannelID {
			continue
		}
		channel := validChannels[random.Intn(len(validChannels))]
		err := session.GuildMemberMove(i.GuildID, userVoiceState.UserID, &channel)
		var logmsg string
		if err != nil {
			logmsg = "Couldn't move %s to %s... oh well..."
		} else {
			logmsg = "Moving %s to %s"
		}
		log.Printf(logmsg, userVoiceState.UserID, channel)
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Fly! You fools!",
		},
	})
}

func momsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.Member
	if member == nil {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Shhhh I'm sleeping.",
			},
		})
		return
	}
	userID := member.User.ID
	guildID := i.GuildID

	callerVoiceState, err := session.State.VoiceState(guildID, userID)
	if err != nil {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> thought they were real slick calling this from a non-voice channel", userID),
			},
		})
		return
	}

	guild, err := session.State.Guild(guildID)
	if err != nil {
		panic(err)
	}

	if guild.AfkChannelID == "" {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No AFK channel for %s", guildID),
			},
		})
		return
	}

	// Check if the caller at least has permission to move members to AFK.
	if !hasMovePermission(session, userID, guild.AfkChannelID) {
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> can't move people to <#%s>", userID, guild.AfkChannelID),
			},
		})
		return
	}

	for _, userVoiceState := range guild.VoiceStates {
		if userVoiceState.ChannelID != callerVoiceState.ChannelID {
			continue
		}
		err := session.GuildMemberMove(guildID, userVoiceState.UserID, &guild.AfkChannelID)
		var logmsg string
		if err != nil {
			logmsg = "Couldn't move %s to %s... oh well..."
		} else {
			logmsg = "Moving %s to %s"
		}
		log.Printf(logmsg, userVoiceState.UserID, guild.AfkChannelID)
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "HIDE! MOM'S HOME!",
		},
	})
}
