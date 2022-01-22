package commands

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "cat",
			Description: "Get yo-self a cat",
		},
		{
			Name:        "dog",
			Description: "Get yo-self a dog",
		},
		{
			Name:        "scatter",
			Description: "Fly! You fools!",
		},
		{
			Name:        "moms-home",
			Description: "HIDE! MOM'S HOME!",
		},
	}
	CommandHandlers = map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate){
		"cat":       catsHandler,
		"dog":       dogsHandler,
		"scatter":   scatterHandler,
		"moms-home": momsHandler,
		"followups": func(session *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   1 << 6,
					Content: "Surprise!",
				},
			})
			msg, err := session.FollowupMessageCreate(session.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				session.FollowupMessageCreate(session.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			session.FollowupMessageEdit(session.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
			})

			time.Sleep(time.Second * 10)

			session.FollowupMessageDelete(session.State.User.ID, i.Interaction, msg.ID)

			session.FollowupMessageCreate(session.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
	}
)
