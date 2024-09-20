package events

import (
	"log"
	"strings"

	"github.com/alexgeraldo/discord-bot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	goodbyeChatID  = config.GetEnv("goodbyechat", "947707167375515669")
	goodbyeMessage = "{} deixou-nos... :sob:"
)

func OnLeaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	// Login member leave event
	log.Printf("%s left the %s guild \n", m.User.Username, s.State.Application.GuildID)

	// Build the message to send
	message := strings.Replace(goodbyeMessage, "{}", m.Mention(), 1)

	// Send goodbye message
	_, err := s.ChannelMessageSend(goodbyeChatID, message)

	// Handler the case where an error occurs while notifying
	if err != nil {
		log.Fatal(err)
	}

}
