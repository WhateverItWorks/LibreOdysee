package proxy

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ProxyImage(c *fiber.Ctx) error {
	imgUrl := c.Query("url")
	hash := c.Query("hash")
	if hash == "" || imgUrl == "" {
		_, err := c.Status(400).WriteString("no hash or url")
		return err
	}

	urlBytes, err := base64.URLEncoding.DecodeString(imgUrl)
	if err != nil {
		return err
	}
	imgUrl = string(urlBytes)

	if !utils.VerifyHMAC(c.Query("url"), hash) {
		_, err := c.Status(400).WriteString("invalid hash")
		return err
	}

	width := "0"
	if c.Query("w") != "" {
		width = c.Query("w")
	}
	height := "0"
	if c.Query("h") != "" {
		height = c.Query("h")
	}

	optionsHash := ""
	if viper.GetBool("IMAGE_CACHE") {
		hasher := sha256.New()
		hasher.Write([]byte(imgUrl + hash + width + height))
		optionsHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		image, err := os.ReadFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash)
		if err == nil {
			_, err := c.Write(image)
			return err
		}
	}

	c.Set("Cache-Control", "public,max-age=31557600")

	client := utils.NewClient()

	requestUrl := "https://thumbnails.odycdn.com/optimize/s:" + width + ":" + height + "/quality:85/plain/" + imgUrl
	if strings.Contains(imgUrl, "static.odycdn.com/emoticons") {
		requestUrl = imgUrl
	}
	
	res, err := client.Get(requestUrl)
	if err != nil {
		return err
	}

	c.Set("Content-Type", res.Header.Get("Content-Type"))

	contentLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	if viper.GetBool("IMAGE_CACHE") && res.StatusCode == 200 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = os.WriteFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash, data, 0644)
		if err != nil {
			return err
		}

		_, err = c.Write(data)
		return err
	} else {
		return c.SendStream(res.Body, contentLen)
	}
}
