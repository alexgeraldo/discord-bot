package events

import (
	"log"
	"strings"

	"github.com/alexgeraldo/discord-bot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	welcomeChatID  = config.GetEnv("welcomechat", "947706647718006794")
	welcomeMessage = "Olá {}! Bem vindo à Comunidade #ACI. Pedimos-te apenas que respeites as <#947706647718006794>!"
)

func OnJoinHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Login member join event
	log.Printf("%s joined the %s guild \n", m.User.Username, s.State.Application.GuildID)

	// Build the message to send
	message := strings.Replace(welcomeMessage, "{}", m.Mention(), 1)

	// Send welcome message
	_, err := s.ChannelMessageSend(welcomeChatID, message)

	// Handler the case where an error occurs while notifying
	if err != nil {
		log.Fatal(err)
	}

}
