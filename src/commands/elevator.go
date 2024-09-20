package commands

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ElevatorCommand = &discordgo.ApplicationCommand{
	Name:        "carousel",
	Description: "Offer a ride on the carousel to a member",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to send on a ride",
			Required:    true, // The parameter is mandatory
		}, {
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "seconds",
			Description: "The duration in seconds of the ride (5-20)",
			Required:    true, // The parameter is mandatory
		},
	},
}

func getUserVoiceState(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceState, error) {
	// Get guild state
	guild, err := s.State.Guild(guildID)
	if err != nil {
		log.Printf("%s", err)
		return nil, fmt.Errorf("Unable to find guild %s", guildID)
	}

	// Look for voice state in guild voice states
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}
	return nil, fmt.Errorf("User <@%s> is not in a voice channel", userID)
}

func canUseElevator(s *discordgo.Session, guildID, commandUserID, targetUserID string) (bool, error) {
	// Fetch guild member for both command user and target user
	commandUser, err := s.GuildMember(guildID, commandUserID)
	if err != nil {
		return false, fmt.Errorf("Failed to fetch command user: %v", err)
	}

	targetUser, err := s.GuildMember(guildID, targetUserID)
	if err != nil {
		return false, fmt.Errorf("Failed to fetch target user: %v", err)
	}

	// Fetch guild to get the roles and their positions
	guild, err := s.Guild(guildID)
	if err != nil {
		return false, fmt.Errorf("Failed to fetch guild: %v", err)
	}

	// Find the highest role position for both users
	highestCommandUserRolePosition := getHighestRolePosition(commandUser, guild)
	highestTargetUserRolePosition := getHighestRolePosition(targetUser, guild)

	// Compare role hierarchy: commandUser must have a higher role than targetUser
	if highestCommandUserRolePosition > highestTargetUserRolePosition {
		return true, nil
	}

	return false, nil
}

func getHighestRolePosition(member *discordgo.Member, guild *discordgo.Guild) int {
	highestPosition := -1

	// Iterate through the member's roles and find the highest position
	for _, memberRoleID := range member.Roles {
		for _, role := range guild.Roles {
			if role.ID == memberRoleID && role.Position > highestPosition {
				highestPosition = role.Position
			}
		}
	}

	return highestPosition
}

// Helper function to filter voice channels and exclude the current channel
func filterVoiceChannels(channels []*discordgo.Channel, excludeChannelID string) []*discordgo.Channel {
	var voiceChannels []*discordgo.Channel
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildVoice && ch.ID != excludeChannelID {
			voiceChannels = append(voiceChannels, ch)
		}
	}
	return voiceChannels
}

func ElevatorHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get command arguments/options
	options := i.ApplicationCommandData().Options
	targetUser := options[0].UserValue(s)
	commandUser := i.Member.User
	durationSeconds := options[1].IntValue()

	// Validate duration between 5 and 20 seconds
	if durationSeconds < 5 || durationSeconds > 20 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The duration must be between 5 and 20 seconds (inclusive).",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Check if the target user is in a voice channel
	voiceState, err := getUserVoiceState(s, i.GuildID, targetUser.ID)
	if err != nil {
		fmt.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The %s must be in a voice channel.", targetUser.GlobalName),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Check if the command user can use the elevator on the target user
	canUse, err := canUseElevator(s, i.GuildID, commandUser.ID, targetUser.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to check permissions: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if !canUse {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use the elevator on this user.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Get a list of all voice channels in the guild
	channels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		log.Printf("An error happened retrieving the channels: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "An error happened, contact an administrator.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	voiceChannels := filterVoiceChannels(channels, voiceState.ChannelID)

	// Move the user to random channels
	go func() {
		for start := time.Now(); time.Since(start).Seconds() < float64(durationSeconds); {
			// Select a random channel
			randomChannel := voiceChannels[rand.Intn(len(voiceChannels))]

			// Move the user to the selected channel
			err = s.GuildMemberMove(i.GuildID, targetUser.ID, &randomChannel.ID)
			if err != nil {
				return
			}

			// Wait for 1 seconds before moving to another channel
			time.Sleep(1 * time.Second)
		}

		// Return the user to the original voice channel
		err = s.GuildMemberMove(i.GuildID, targetUser.ID, &voiceState.ChannelID)
		if err != nil {
			return
		}
	}()

	// Respond to acknowledge the command
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s vai dar uma voltinha de %d segundos no carrossel.", targetUser.Mention(), durationSeconds),
		},
	})
}
