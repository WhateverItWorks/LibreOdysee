package api

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/dustin/go-humanize"
)

type Live struct {
	ClaimId      string
	RelTime      string
	Time         string
	ThumbnailUrl string
	StreamUrl    string
	Live         bool
}

func GetLive(claimId string) (Live, error) {
	data, err := utils.RequestJSON("https://api.odysee.live/livestream/is_live?channel_claim_id="+claimId, nil)
	if err != nil {
		return Live{}, err
	}

	if !data.Get("success").Bool() {
		return Live{}, fmt.Errorf(data.Get("error").String())
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", data.Get("data.Start").String())
	if err != nil {
		return Live{}, err
	}

	thumbnail := data.Get("data.ThumbnailURL").String()
	thumbnail = base64.URLEncoding.EncodeToString([]byte(thumbnail))
	thumbnail = "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail)

	streamUrl := strings.ReplaceAll(data.Get("data.VideoURL").String(), "https://cloud.odysee.live", "/live")
	streamUrl = strings.ReplaceAll(streamUrl, "https://cdn.odysee.live", "/live")

	return Live{
		RelTime:      humanize.Time(timestamp),
		Time:         timestamp.Format("Jan 2, 2006 03:04 PM"),
		ThumbnailUrl: thumbnail,
		StreamUrl:    streamUrl,
		Live:         data.Get("data.Live").Bool(),
	}, nil
}
