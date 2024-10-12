package main

import (
	"fmt"
	"github.com/alexgeraldo/discord-bot/types"
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
	animeChatID    = config.GetEnv("animechat", "1030120857030361149")
	newsChatID     = config.GetEnv("newschat", "1286744761088086157")
)

var registeredCommands = make(map[string]types.CommandInfo)

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
		if command, ok := registeredCommands[commandName]; ok {
			command.Handler(s, i)

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

	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

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
	episodesLastChecked := time.Now()
	now := time.Now()
	newsLastChecked := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) //time.Now()

	// Add tasks to cron scheduler
	log.Println("Adding cron tasks...")

	tasks.CrunchyrollEpisodesNotification(s, animeChatID, &episodesLastChecked)
	_, err = c.AddFunc("@every 5m", func() {
		tasks.CrunchyrollEpisodesNotification(s, animeChatID, &episodesLastChecked)
	})
	if err != nil {
		log.Fatalf("Error adding crunchyroll episodes to cron tasks: %v", err)
	}

	tasks.CrunchyrollArticlesNotification(s, newsChatID, &newsLastChecked)
	_, err = c.AddFunc("@every 30m", func() {
		tasks.CrunchyrollArticlesNotification(s, newsChatID, &newsLastChecked)
	})
	if err != nil {
		log.Fatalf("Error adding crunchyroll articles to cron tasks: %v", err)
	}

	// Register slash commands to the bot
	log.Println("Adding commands...")

	err = commands.RegisterHelloCommand(s, guildID, registeredCommands)
	if err != nil {
		log.Println(err)
	}
	err = commands.RegisterRoastCommand(s, guildID, registeredCommands)
	if err != nil {
		log.Println(err)
	}
	err = commands.RegisterElevatorCommand(s, guildID, registeredCommands)
	if err != nil {
		log.Println(err)
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
		fmt.Println(registeredCommands)
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.Command.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Command.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
