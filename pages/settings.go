package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func SettingsHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=604800")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self'; connect-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	theme := viper.GetString("DEFAULT_SETTINGS.theme")
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}
	
	return c.Render("settings", fiber.Map{
		"config": viper.AllSettings(),
		"theme": theme,
	})
}

func DefaultSettingsHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=604800")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Access-Control-Allow-Origin", "*")
		
	if viper.GetString("DEFAULT_SETTINGS.theme") == "" {
		return c.JSON(fiber.Map{
		"theme": "system",
		"relatedVideos": true,
		"nsfw": false,
		"autoplay": false,
		"speed": "1",
		"quality": "0",
		"commentWarning": true,
		"sponsorblock": fiber.Map{
			"sponsor": true,
			"selfpromo": true,
			"interaction": true,
			"intro": false,
			"outro": false,
			"preview": false,
			"filler": false,
		},
	})
	}

	viper.SetDefault("DEFAULT_SETTINGS.commentWarning", true)
		
	return c.JSON(fiber.Map{
		"theme": viper.GetString("DEFAULT_SETTINGS.theme"),
		"relatedVideos": viper.GetBool("DEFAULT_SETTINGS.relatedVideos"),
		"nsfw": viper.GetBool("DEFAULT_SETTINGS.nsfw"),
		"autoplay": viper.GetBool("DEFAULT_SETTINGS.autoplay"),
		"speed": viper.GetString("DEFAULT_SETTINGS.speed"),
		"quality": viper.GetString("DEFAULT_SETTINGS.quality"),
		"commentWarning": viper.GetBool("DEFAULT_SETTINGS.commentWarning"),
		"sponsorblock": fiber.Map{
			"sponsor": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.sponsor"),
			"selfpromo": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.selfpromo"),
			"interaction": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.interaction"),
			"intro": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.intro"),
			"outro": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.outro"),
			"preview": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.preview"),
			"filler": viper.GetBool("DEFAULT_SETTINGS.sponsorblock.filler"),
		},
	})
}