package tasks

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	crunchyEpisodesEndpoint = "https://feeds.feedburner.com/crunchyroll/rss/anime"
	crunchyNewsEndpoint     = "https://cr-news-api-service.prd.crunchyrollsvc.com/v1/pt-BR/rss"
	thumbnailList           = []string{
		"https://cdn.discordapp.com/attachments/808470228378845275/1286477872747380816/208942.512.webp?ex=66ee0d62&is=66ecbbe2&hm=3e10b8d6df5e68fb75576898c1721f7bca7d011ca591de2b198eddd977bba191&",
		"https://cdn.discordapp.com/attachments/808470228378845275/1286480363547525180/Bond_Forger_Anime.webp?ex=66ee0fb4&is=66ecbe34&hm=e3ca84fcf66920cc9fae73a4d06aa20229ec8e13bf036065b18b76feea8fda05&",
		"https://cdn.discordapp.com/attachments/808470228378845275/1286480363798925363/07193dd7b9182d275deb3b0b789e0588.png?ex=66ee0fb4&is=66ecbe34&hm=52e0180ac6e098919ee07616689fe882a3f88bc31a343287647ef6a9ce7085fd&",
		"https://cdn.discordapp.com/attachments/808470228378845275/1286481919546232842/1656860964595.webp?ex=66ee1127&is=66ecbfa7&hm=09df550e905a152da8b552b81d661fc05ce8d90f2ac3366fc1263d9103d3f085&",
		"https://cdn.discordapp.com/attachments/808470228378845275/1286481919868932187/subir.png?ex=66ee1127&is=66ecbfa7&hm=b6a1581a95605563974d3f182f02bf160acd9c6e15ad1c3fb2dec9d6776aa669&",
	}
)

func removeHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}

func CrunchyrollEpisodesNotification(s *discordgo.Session, channelID string, lastCheck *time.Time) {
	log.Println("Checking for new episodes in the Crunchyroll RSS feed...")

	// Parse from rss feed url
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(crunchyEpisodesEndpoint)

	// Parse anime rss feed
	if err != nil {
		log.Printf("An error occurred parsing: %s for %s\n", err, crunchyEpisodesEndpoint)
		return
	}

	for _, item := range feed.Items {
		// Check if the episode was published after last check
		pubTime, err := time.Parse(time.RFC1123, item.Published)
		if err != nil {
			log.Printf("An error occurred while checking date: %v\n", err)
			continue
		}

		// Build message for new and valid item (is not a dub and is new)
		if !strings.Contains(item.Title, " Dub)") && pubTime.After(*lastCheck) {
			log.Printf("New episode found '%s'.\n", item.Title)

			// Get thumbnail
			thumbnailURL := thumbnailList[rand.Intn(len(thumbnailList))] // Default funny images to replace in case of error

			if mediaThumbnails, ok := item.Extensions["media"]["thumbnail"]; ok {
				thumbnailURL = mediaThumbnails[0].Attrs["url"]
			} else {
				log.Printf("Erro ao obter a cover do episodio: %v", err)
			}

			// Create the embed with title and description for the episode launched
			embed := &discordgo.MessageEmbed{
				Title:       item.Title,
				URL:         item.GUID,
				Description: fmt.Sprintf("%s\n\n**Lançado em:** %s", removeHTMLTags(item.Description), pubTime.Format("02 Jan 2006 15:04")),
				Color:       0x00FF00, // Green color
				Timestamp:   time.Now().Format(time.RFC3339),
				Image: &discordgo.MessageEmbedImage{
					URL: thumbnailURL,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Fonte: Crunchyroll",
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

func CrunchyrollArticlesNotification(s *discordgo.Session, channelID string, lastCheck *time.Time) {
	log.Println("Checking for new articles in the Crunchyroll RSS feed...")

	// Parse from rss feed url
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(crunchyNewsEndpoint)

	// Parse anime rss feed
	if err != nil {
		log.Printf("An error occurred parsing: %s for %s\n", err, crunchyNewsEndpoint)
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
				Timestamp:   pubTime.Format(time.RFC3339), // Use the publishing time for timestamp
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

	// Update last check timestamp
	*lastCheck = time.Now()
}
