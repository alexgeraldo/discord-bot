package commands

import (
	"fmt"
	"github.com/alexgeraldo/discord-bot/types"
	"log"

	"github.com/bwmarrin/discordgo"
)

var helloCommand = &discordgo.ApplicationCommand{
	Name:        "hello",
	Description: "Completes for 'Hello World!' message",
}

func helloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Answer with 'World!' to complete user message
	log.Printf("Completing %v for 'Hello World!' message\n", i.Member.User.Username)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "World!",
		},
	})

	// Handle error
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterHelloCommand(s *discordgo.Session, guildID string, registeredMap map[string]types.CommandInfo) error {
	// Create Application Command
	helloCommand, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, helloCommand)
	if err != nil {
		return fmt.Errorf("error creating hello command: %v", err)
	}

	// Add Application Info to the registeredMap
	command := types.CommandInfo{
		Command: helloCommand,
		Handler: helloHandler,
	}
	registeredMap[helloCommand.Name] = command

	// Successful nil
	return nil
}
