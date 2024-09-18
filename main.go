package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// Bot commands
var (
	commandsList = []*discordgo.ApplicationCommand{
		commands.helloCommand,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
	}
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it is interaction
	if i.Type == discordgo.InteractionApplicationCommand {

		// Get the name of the slash command
		commandName := i.ApplicationCommandData().Name

		// Execute the handler for the slash command
		if handler, ok := commandHandlers[commandName]; ok {
			handler(s, i)

			// Warn the admin that the handler is not implemented
		} else {
			fmt.Printf("Handler for '%v' slash-command is not implemented!", commandName)

			// Respond the user with generic error message
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "An error had occurred! Contact thebot administrator.",
				},
			})

			// Handler the case where an error occurs while responding to interaction
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	// Get discord bot token
	TOKEN := "TOKEN HERE"

	// Create discord bot session
	s, err := discordgo.New("Bot " + TOKEN)

	// Check if an error happened creating session
	if err != nil {
		log.Fatalf("Error starting bot session: %v", err)
	}

	// Add interaction handler to the bot
	s.AddHandler(interactionHandler)

	// Start the discord bot session
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Register slash commands to the bot
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	// Hold the thread with a channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
