package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

var (
	animescheduleEpisodesEndpoint = "https://animeschedule.net/subrss.xml"
	// From documentation https://animeschedule.net/api/v3/documentation
	apiBaseUrl   = "https://animeschedule.net/api/v3/anime"
	imageBaseUrl = "https://img.animeschedule.net/production/assets/public/img"
)

// Function to get the cover of an anime from the API using the slug extracted from the feed link
func getAnimeCover(slug string) (string, error) {
	// Make the GET request to the API using the slug from the link
	apiURL := fmt.Sprintf("%s/%s", apiBaseUrl, slug)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error retrieving anime cover: %v from %s", err, apiURL)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %v from %s", err, apiURL)
	}

	// Structure to deserialize the JSON
	var result struct {
		ImageVersionRoute string `json:"imageVersionRoute"`
	}

	// Deserialize the JSON
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Construct the full URL for the anime image
	imageURL := fmt.Sprintf("%s/%s", imageBaseUrl, result.ImageVersionRoute)
	return imageURL, nil
}

func AnimescheduleEpisodesNotification(s *discordgo.Session, channelID string, lastCheck *time.Time) {
	// Parse from rss feed url
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(animescheduleEpisodesEndpoint)

	// Parse anime rss feed
	if err != nil {
		log.Printf("An error occurred parsing: %s for %s\n", err, animescheduleEpisodesEndpoint)
		return
	}

	for _, item := range feed.Items {
		// Check if the episode was published after last check
		pubTime, err := time.Parse(time.RFC1123, item.Published)
		if err != nil {
			log.Printf("An error occurred while checking date: %v\n", err)
			continue
		}

		if pubTime.After(*lastCheck) {
			log.Printf("New episode found '%s'.\n", item.GUID)

			// Extract anime slug from the link
			linkParts := strings.Split(item.Link, "/")
			slug := linkParts[len(linkParts)-1]

			// Obter a cover do anime via API
			imageURL, err := getAnimeCover(slug)
			if err != nil {
				log.Printf("Erro ao obter a cover do anime: %v", err)
				imageURL = thumbnailList[rand.Intn(len(thumbnailList))] // Default funny images to replace in case of error
			}

			// Create the embed with title and description for the episode lauched
			embed := &discordgo.MessageEmbed{
				Title:       item.Title,
				URL:         item.Link,
				Description: fmt.Sprintf("**%s** foi lançado com legendas!\n\n**Publicado em:** %s", item.GUID, pubTime.Format("02 Jan 2006 15:04")),
				Color:       0x00FF00, // Green color
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Link para o Episódio",
						Value:  item.Link,
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
				Image: &discordgo.MessageEmbedImage{
					URL: imageURL,
				},
			}

			// Send the embed
			_, err = s.ChannelMessageSendEmbed(channelID, embed)
			if err != nil {
				log.Printf("Error sending the emebed to Discord: %v\n", err)
			}
		}
	}

	// Update last check timestamp
	*lastCheck = time.Now()
}
