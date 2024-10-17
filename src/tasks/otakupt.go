package tasks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	"log"
	"regexp"
	"time"
)

var otakuNewsEndpoint = "https://www.otakupt.com/category/anime/feed/"

func extractImageURLFromDescription(item *gofeed.Item) (string, bool) {
	// Check if the item description contains an <img> tag with src attribute
	re := regexp.MustCompile(`<img[^>]*src="([^"]+)"[^>]*>`)
	match := re.FindStringSubmatch(item.Description)

	if len(match) > 1 {
		// Return the URL and true if found
		return match[1], true
	}
	// Return an empty string and false if no image URL is found
	return "", false
}

func OtakuArticlesNotification(s *discordgo.Session, channelID string, lastCheck *time.Time) {
	log.Println("Checking for new articles in the OtakuPT RSS feed...")

	// Parse from rss feed url
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(otakuNewsEndpoint)

	// Parse anime rss feed
	if err != nil {
		log.Printf("An error occurred parsing: %s for %s\n", err, otakuNewsEndpoint)
		return
	}

	// Go through the articles and checks what is new
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
				Description: removeHTMLTags(item.Description),
				URL:         item.Link,
				Color:       0x00FF00,                     // Green color
				Timestamp:   pubTime.Format(time.RFC3339), // Use the publishing time for timestamp
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Fonte: OtakuPT",
				},
			}

			// Extract media:thumbnail from custom fields
			if mediaThumbnail, ok := extractImageURLFromDescription(item); ok {
				if len(mediaThumbnail) > 0 {
					embed.Image = &discordgo.MessageEmbedImage{
						URL: mediaThumbnail,
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

	// Update last check timestamp
	*lastCheck = time.Now()
}
