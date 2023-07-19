package pages

import (
	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"codeberg.org/librarian/feeds"
	"github.com/spf13/viper"
)

func ChannelRSSHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=1800")
	c.Set("Content-Type", "application/rss+xml")

	channel, err := api.GetChannel(c.Params("channel"))
	if err != nil {
		return err
	}
	if channel.Id == "" {
		c.Set("Content-Type", "text/plain")
		_, err := c.Status(404).WriteString("404 Not Found\nERROR: Unable to find channel")
		return err
	}
	claims, err := channel.GetClaims(1)
	if err != nil {
		return err
	}

	image, err := utils.UrlEncode(viper.GetString("DOMAIN") + channel.Thumbnail)
	if err != nil {
		_, err := c.Status(500).WriteString("500 Internal Server Error\nERROR: " + err.Error())
		return err
	}

	feed := &feeds.Feed{
		Title:       channel.Name + " - Librarian",
		Link:        &feeds.Link{Href: channel.Url},
		Image:       &feeds.Image{Url: image},
		Description: channel.DescriptionTxt,
	}

	feed.Items = []*feeds.Item{}

	for i := 0; i < len(claims); i++ {
		item := &feeds.Item{
			Title:       claims[i].Title,
			Link:        &feeds.Link{Href: claims[i].Url},
			Description: "<img width=\"480\" src=\"" + viper.GetString("DOMAIN") + claims[i].ThumbnailUrl + "\"><br><br>" + string(claims[i].Description),
			Created:     claims[i].Time,
			Enclosure:   &feeds.Enclosure{},
		}

		if c.Query("odyseeLink") == "true" {
			item.Link.Href = claims[i].OdyseeUrl
		}

		if c.Query("enclosure") == "true" {
			stream, err := api.GetStream(claims[i].LbryUrl)
			if err != nil {
				return err
			}
			url, err := utils.UrlEncode(stream.URL)
			if err != nil {
				_, err := c.Status(500).WriteString("500 Internal Server Error\nERROR: " + err.Error())
				return err
			}
			item.Enclosure.Url = url
			item.Enclosure.Type = claims[i].MediaType
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	_, err = c.Write([]byte(rss))
	return err
}
