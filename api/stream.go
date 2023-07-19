package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"codeberg.org/librarian/librarian/utils"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var streamCache = cache.New(30*time.Minute, 15*time.Minute)

type Stream struct {
	Type        string
	URL         string
	FallbackURL string
	HLS         bool
}

func GetStream(video string) (Stream, error) {
	cacheData, found := streamCache.Get(video + "-stream")
	if found {
		return cacheData.(Stream), nil
	}

	reqDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "get",
		"params": map[string]interface{}{
			"uri":       video,
			"save_file": false,
		},
		"id": time.Now().Unix(),
	}
	reqData, err := json.Marshal(reqDataMap)
	if err != nil {
		return Stream{}, err
	}

	data, err := utils.RequestJSON(viper.GetString("STREAMING_API_URL")+"?m=get", bytes.NewBuffer(reqData))
	if err != nil {
		return Stream{}, err
	}

	if data.Get("error.message").String() != "" {
		return Stream{}, err 
	}

	streamUrl := data.Get("result.streaming_url").String()
	streamUrl = strings.ReplaceAll(streamUrl, "source.odycdn.com", "player.odycdn.com")
	if viper.GetString("VIDEO_STREAMING_URL") != "" {
		streamUrl = strings.ReplaceAll(streamUrl, "http://localhost:5280", viper.GetString("VIDEO_STREAMING_URL"))
		streamUrl = strings.ReplaceAll(streamUrl, "https://player.odycdn.com", viper.GetString("VIDEO_STREAMING_URL"))
	}

	stream, err := checkStream(streamUrl)
	if err != nil {
		return Stream{}, err
	}

	streamCache.Set(video+"-stream", stream, cache.DefaultExpiration)
	return stream, nil
}

func checkStream(url string) (Stream, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return Stream{}, err
	}

	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "https://odysee.com")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/109.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Stream{}, err
	}

	if res.StatusCode == 403 {
		return Stream{}, fmt.Errorf("this content cannot be accessed due to a DMCA request")
	}
	
	return Stream{
		Type:        res.Header.Get("Content-Type"),
		URL:         res.Request.URL.String(),
		FallbackURL: url,
		HLS: res.Header.Get("Content-Type") == "application/x-mpegurl",
	}, nil
}
