package tasks

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

var (
	newsRssSource = "https://cr-news-api-service.prd.crunchyrollsvc.com/v1/pt-BR/rss"
)

func NotifyNewArticle(s *discordgo.Session, channelID string, lastCheck *time.Time) {
	// Parse from rss feed url
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(newsRssSource)

	// Parse anime rss feed
	if err != nil {
		log.Printf("An error occurred parsing: %s for %s\n", err, rssSource)
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
			// Create an embed with the article details
			embed := &discordgo.MessageEmbed{
				Title:       item.Title,
				Description: item.Description,
				URL:         item.Link,
				Color:       0x00FF00,                     // Green color
				Timestamp:   pubTime.Format(time.RFC3339), // Use the publish time for timestamp
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Fonte: Crunchyroll",
				},
			}

			// Extract media:thumbnail from custom fields
			if mediaThumbnail, ok := item.Extensions["media"]["thumbnail"]; ok {
				if len(mediaThumbnail) > 0 {
					embed.Image = &discordgo.MessageEmbedImage{
						URL: mediaThumbnail[0].Attrs["url"],
					}
				}
			}

			// Send the embed to the specified channel
			_, err := s.ChannelMessageSendEmbed(channelID, embed)
			if err != nil {
				log.Printf("Error sending the emebed to Discord: %v\n", err)
			}
		}
	}
}
