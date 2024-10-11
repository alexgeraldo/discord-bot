package types

import "github.com/bwmarrin/discordgo"

// CommandInfo holds the command definition and its handler
type CommandInfo struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
