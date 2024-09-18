package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var helloCommand = &discordgo.ApplicationCommand{
	Name:        "hello",
	Description: "Completes for 'Hello World!' message",
}

func helloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Answer with 'World!" to complete user message
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
