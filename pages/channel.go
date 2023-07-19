package pages

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/WhateverItWorks/LibreOdysee/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ChannelHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=1800")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	theme := viper.GetString("DEFAULT_SETTINGS.theme")
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}
	
	page := 1
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channel, err := api.GetChannel(c.Params("channel"))
	if err != nil {
		return err
	}
	channel.GetFollowers()

	if channel.Id == "" {
		return c.Status(404).Render("errors/notFound", fiber.Map{})
	}

	if channel.ValueType != "channel" {
		return ClaimHandler(c)
	}

	claims, err := channel.GetClaims(page)
	if err != nil {
		return err
	}
	sort.Slice(claims, func(i int, j int) bool {
		return claims[i].Timestamp > claims[j].Timestamp
	})

	return c.Render("channel", fiber.Map{
		"channel": channel,
		"config":  viper.AllSettings(),
		"claims":  claims,
		"theme":   theme,
		"query": fiber.Map{
			"page":        fmt.Sprint(page),
			"prevPageIs0": (page - 1) == 0,
			"nextPage":    fmt.Sprint(page + 1),
			"prevPage":    fmt.Sprint(page - 1),
		},
	})
}

func ChannelApiHandler(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET")

	page := 0
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channel, err := api.GetChannel(c.Params("channel"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"status": "500",
		})
	}
	channel.GetFollowers()

	if channel.Id == "" {
		return c.Status(404).JSON(fiber.Map{
			"error": "Channel not found",
			"status": "404",
		})
	}

	if channel.ValueType != "channel" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid value type, excepted channel",
			"status": "400",
		})
	}

	if page != 0 {
		claims, err := channel.GetClaims(page)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
				"status": "500",
			})
		}
		sort.Slice(claims, func(i int, j int) bool {
			return claims[i].Timestamp > claims[j].Timestamp
		})
		channel.Claims = claims
	}

	return c.JSON(fiber.Map{
		"channel": channel,
		"status": "200",
	})
}
