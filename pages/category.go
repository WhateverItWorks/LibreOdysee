package pages

import (
	"github.com/WhateverItWorks/LibreOdysee/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func CategoryHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=1800")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	categories, err := api.GetCategoryData()
	if err != nil {
		return err
	}
	categoriesArr, err := api.GetOrderedCategoriesArray()
	if err != nil {
		return err
	}

	categoryName := "featured"
	if c.Params("category") != "" {
		categoryName = c.Params("category")
	}

	claims, err := categories[categoryName].GetCategoryClaims(1, c.Cookies("nsfw") == "true")
	if err != nil {
		return err
	}
	
	theme := viper.GetString("DEFAULT_SETTINGS.theme")
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}

	return c.Render("category", fiber.Map{
		"config":     viper.AllSettings(),
		"category":   categories[categoryName],
		"categories": categoriesArr,
		"claims":     claims,
		"theme":      theme,
	})
}

func CategoryApiHandler(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET")

	categories, err := api.GetCategoryData()
	if err != nil {
		return err
	}

	categoryName := "featured"
	if c.Params("category") != "" {
		categoryName = c.Params("category")
	}

	claims, err := categories[categoryName].GetCategoryClaims(1, c.Query("nsfw") == "true")
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"category":   categories[categoryName],
		"claims":     claims,
		"status": "200",
	})
}
