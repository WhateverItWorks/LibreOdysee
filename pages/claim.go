package pages

import (
	"net/url"
	"strings"

	"github.com/WhateverItWorks/LibreOdysee/api"
	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ClaimHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=21600")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'self'; script-src blob: 'self'; connect-src *; media-src * data: blob:; block-all-mixed-content")

	theme := viper.GetString("DEFAULT_SETTINGS.theme")
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}
	showRelated := viper.GetString("DEFAULT_SETTINGS.showRelated")
	if c.Cookies("showRelated") != "" {
		showRelated = c.Cookies("showRelated")
	}
	autoplay := viper.GetString("DEFAULT_SETTINGS.autoplay")
	if c.Cookies("autoplay") != "" {
		autoplay = c.Cookies("autoplay")
	}
	nojs := c.Query("nojs") == "1"
	settings := fiber.Map{
		"theme": theme,
		"nojs":  nojs,
		"showRelated": showRelated,
		"autoplay": autoplay,
	}

	claimData, err := api.GetClaim("lbry://" + c.Params("channel") + "/" + c.Params("claim"))
	if err != nil {
		if strings.ContainsAny(err.Error(), "NOT_FOUND") {
			return c.Status(404).Render("errors/notFound", fiber.Map{"theme": theme})
		}
		return err
	}

	if claimData.ValueType == "repost" {
		repostLink, err := utils.LbryTo(claimData.Repost)
		if err != nil {
			return err
		}
		return c.Redirect(repostLink["rel"])
	}

	if utils.Contains(viper.GetStringSlice("blocked_claims"), claimData.Id) {
		return c.Status(451).Render("errors/blocked", fiber.Map{
			"claim": claimData,
			"theme": theme,
		})
	}

	if claimData.HasFee {
		return c.Render("errors/hasFee", fiber.Map{
			"claim": claimData,
			"theme": theme,
		})
	}

	related, err := api.GetRelated(claimData.Title, c.Cookies("nsfw") == "true", claimData.Id)
	if err != nil {
		return err
	}

	if claimData.MediaType == "" && claimData.ValueType == "stream" {
		live, err := api.GetLive(claimData.Channel.Id)
		if err != nil && err.Error() != "no data associated with claim id" {
			return err
		}

		if !viper.GetBool("ENABLE_LIVESTREAM") {
			return c.Render("errors/liveDisabled", fiber.Map{
				"switchUrl": c.Path(),
				"settings":  settings,
			})
		}

		return c.Render("live", fiber.Map{
			"live":     live,
			"claim":    claimData,
			"settings": settings,
			"config":   viper.AllSettings(),
		})
	}

	stream, err := api.GetStream(claimData.LbryUrl)
	if err != nil {
		if err.Error() == "this content cannot be accessed due to a DMCA request" {
			return c.Status(451).Render("errors/dmca", nil)
		}
		return err
	}

	if viper.GetBool("ENABLE_STREAM_PROXY") && claimData.StreamType != "document" {
		strUrl, err := url.Parse(stream.URL)
		if err != nil {
			return err
		}
		stream.URL = "/stream" + strUrl.Path
		stream.FallbackURL = strings.ReplaceAll(stream.FallbackURL, "https://player.odycdn.com", "/stream")
	}

	comments := api.Comments{}
	if nojs {
		comments, err = claimData.GetComments("", 3, 25, 1)
		if err != nil {
			return err
		}
	}

	switch claimData.StreamType {
	case "document":
		props := fiber.Map{
			"stream":   stream,
			"claim":    claimData,
			"comments": comments,
			"settings": settings,
			"config":   viper.AllSettings(),
		}

		body := []byte{}
		if strings.HasPrefix(stream.Type, "text") {
			body, err = utils.Request(stream.URL, 500000, utils.Data{Bytes: nil})
			if err != nil {
				if strings.ContainsAny(err.Error(), "over byte limit") {
				 props["download"] = true
				 return c.Render("claim", props)
			 }
			 return err
		 }
		}
		
		switch stream.Type {
		case "text/html":
			props["document"] = utils.ProcessDocument(string(body), false)
		case "text/plain":
			props["document"] = string(body)
		case "text/markdown":
			props["document"] = utils.ProcessDocument(string(body), true)
		case "application/pdf":
			c.Set("Content-Security-Policy", "default-src 'self'; script-src blob: 'self'; connect-src *; frame-src 'self' https://player.odycdn.com; block-all-mixed-content")
			props["document"] = `<iframe class="pdf" src="` + stream.URL + `" width="100%"></iframe>`
		default:
			props["download"] = true
			return c.Render("claim", props)
		}

		return c.Render("claim", props)
	case "video":
		if stream.HLS {
			c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; img-src *; script-src blob: 'self'; connect-src *; media-src * data: blob:; block-all-mixed-content")
		}

		return c.Render("claim", fiber.Map{
			"stream":      stream,
			"claim":       claimData,
			"relatedVids": related,
			"comments":    comments,
			"settings":    settings,
			"config":      viper.AllSettings(),
		})
	default:
		return c.Render("claim", fiber.Map{
			"stream":   stream,
			"download": true,
			"comments": comments,
			"claim":    claimData,
			"settings": settings,
			"config":   viper.AllSettings(),
		})
	}
}
