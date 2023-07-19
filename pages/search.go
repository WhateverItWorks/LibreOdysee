package pages

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/WhateverItWorks/LibreOdysee/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func SearchHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "private")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	theme := viper.GetString("DEFAULT_SETTINGS.theme")
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}
	
	page := 1
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}
	pageParam, err = strconv.Atoi(c.FormValue("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	query := c.FormValue("q")
	if query == "" && c.Query("q") != "" {
		query = c.Query("q")
	}

	if len(query) <= 2 {
		return c.Render("search", fiber.Map{
			"results":   nil,
			"lenUnder3": true,
			"theme":     theme,
			"query": fiber.Map{
				"query": query,
			},
		})
	}

	results, err := api.Search(query, page, "file,channel", c.Cookies("nsfw") == "true", 12)
	if err != nil {
		return err
	}
	sort.Slice(results, func(i int, j int) bool {
		valueType := reflect.ValueOf(&results[i]).Elem().Elem().FieldByName("ValueType").String()
		if valueType == "channel" {
			return true
		} else {
			return false
		}
	})

	return c.Render("search", fiber.Map{
		"results": results,
		"theme":   theme,
		"query": fiber.Map{
			"query":       query,
			"page":        fmt.Sprint(page),
			"prevPageIs0": (page - 1) == 0,
			"nextPage":    fmt.Sprint(page + 1),
			"prevPage":    fmt.Sprint(page - 1),
		},
	})
}
