package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alexgeraldo/discord-bot/commands"
	"github.com/alexgeraldo/discord-bot/config"
	"github.com/alexgeraldo/discord-bot/events"
	"github.com/alexgeraldo/discord-bot/tasks"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// Bot parameters
var (
	guildID        = config.GetEnv("guild", "") // GuildID to register commands. If not passed - bot registers commands globally.
	botToken       = config.GetEnv("token", "")
	removeCommands = config.GetEnvAsBool("rmcmd", true)
)

// Bot commands
var (
	commandsList = []*discordgo.ApplicationCommand{
		commands.HelloCommand,
		commands.RoastCommand,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello": commands.HelloHandler,
		"roast": commands.RoastHandler,
	}
)

// Bot events
var (
	eventHandlers = []interface{}{
		interactionHandler,
		events.OnJoinHandler,
		events.OnLeaveHandler,
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

	// Create discord bot session
	s, err := discordgo.New("Bot " + botToken)

	// Check if an error happened creating session
	if err != nil {
		log.Fatalf("Error starting bot session: %v", err)
	}

	// Add bot event handlers
	for _, handler := range eventHandlers {
		s.AddHandler(handler)
	}

	// Start the discord bot session
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Create cron schedule for tasks
	c := cron.New()

	// RSS Feed last checked memory variable
	rssLastChecked := time.Now()

	// Add tasks to cron scheduler
	log.Println("Adding cron tasks...")
	log.Println("Checking for new animes in the RSS feed...")
	tasks.NotifyNewAnime(s, "1030120857030361149", &rssLastChecked)
	c.AddFunc("@every 5m", func() {
		log.Println("Checking for new animes in the RSS feed...")
		tasks.NotifyNewAnime(s, "1030120857030361149", &rssLastChecked)
	})

	// Register slash commands to the bot
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandsList))
	for i, v := range commandsList {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Start cron scheduler
	c.Start()

	defer s.Close()

	// Hold the thread with a channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if removeCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
