package proxy

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleStream(c *fiber.Ctx) error {
	req, err := http.NewRequest("GET", "https://player.odycdn.com" + strings.TrimPrefix(c.Path(), "/stream"), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/109.0")
	req.Header.Set("Origin", "https://odysee.com")
	if c.Get("Range") != "" {
		req.Header.Set("Range", c.Get("Range"))
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Accept-Ranges", "bytes")
	c.Set("Content-Length", res.Header.Get("Content-Length"))
	c.Set("Content-Type", res.Header.Get("Content-Type"))
	if res.Header.Get("Content-Range") != "" {
		c.Set("Content-Range", res.Header.Get("Content-Range"))
	}
	if res.Header.Get("Location") != "" {
		c.Set("Location", strings.ReplaceAll(res.Request.URL.String(), "https://player.odycdn.com", "/stream"))
		res.StatusCode = 308
	}
	c.Status(res.StatusCode)
	
	return c.SendStream(res.Body)
}

func HandleLive(c *fiber.Ctx) error {
	req, err := http.NewRequest("GET", "https://cloud.odysee.live" + strings.TrimPrefix(c.Path(), "/live"), nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
	req.Header.Set("Origin", "https://odysee.com/")
	req.Header.Set("Referer", "https://odysee.com/")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Content-Type", res.Header.Get("Content-Type"))
	c.Status(res.StatusCode)

	if strings.HasSuffix(c.Path(), ".m3u8") {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`(?m)^/live`)
		newBody := re.ReplaceAllString(string(body), "/live/live")
		re2 := regexp.MustCompile(`(?m)^/[0-9]{3}`)
		newBody = re2.ReplaceAllString(newBody, "/live$0")
		newBody = strings.ReplaceAll(newBody, "https://cloud.odysee.live", "/live")
		newBody = strings.ReplaceAll(newBody, "https://cdn.odysee.live", "/live")

		return c.SendString(newBody)
	}
	
	return c.SendStream(res.Body)
}