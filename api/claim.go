package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var claimCache = cache.New(30*time.Minute, 15*time.Minute)

type Claim struct {
	Url          string
	RelUrl       string
	LbryUrl      string
	OdyseeUrl    string
	Id           string
	Channel      Channel
	Title        string
	ThumbnailUrl string
	Description  template.HTML
	License      string
	Views        int64
	Likes        int64
	Dislikes     int64
	Tags         []string
	Timestamp    int64
	RelTime      string
	Date         string
	Time         time.Time
	Duration     string
	MediaType    string
	Repost       string
	ValueType    string
	StreamType   string
	HasFee       bool
}

func GetClaim(lbryUrl string) (Claim, error) {
	cacheData, found := claimCache.Get(lbryUrl)
	if found {
		return cacheData.(Claim), nil
	}

	lbryUrl = strings.TrimSuffix(lbryUrl, "/")

	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     []string{lbryUrl},
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=resolve", bytes.NewBuffer(resolveData))
	if err != nil {
		return Claim{}, err
	}
	data = data.Get("result.lbry*")

	if data.Get("error").String() != "" {
		return Claim{}, fmt.Errorf("API Error: " + data.Get("error").String())
	}

	claim, err := ProcessClaim(data)
	if err != nil {
		return Claim{}, err
	}
	claim.GetViews()
	claim.GetRatings()

	claimCache.Set(lbryUrl, claim, cache.DefaultExpiration)
	return claim, nil
}

func GetClaims(lbryUrls []string, getFollowers bool, getViews bool) ([]interface{}, error) {
	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     lbryUrls,
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=resolve", bytes.NewBuffer(resolveData))
	if err != nil {
		return nil, err
	}

	claims := []interface{}{}
	data.Get("result").ForEach(func(key, value gjson.Result) bool {
		if value.Get("value_type").String() == "channel" {
			channel, _ := ProcessChannel(value)
			if getFollowers {
				channel.GetFollowers()
			}
			claims = append(claims, channel)
		} else {
			claim, _ := ProcessClaim(value)
			if getViews {
				claim.GetViews()
			}
			claims = append(claims, claim)
		}
		return true
	})

	return claims, nil
}

func ProcessClaim(claimData gjson.Result) (Claim, error) {
	if claimData.Get("value_type").String() == "channel" {
		return Claim{}, fmt.Errorf("value type is channel")
	}

	claim := Claim{
		LbryUrl: claimData.Get("canonical_url").String(),
		Id:      claimData.Get("claim_id").String(),
		Channel: Channel{
			Name:        claimData.Get("signing_channel.name").String(),
			Title:       claimData.Get("signing_channel.value.title").String(),
			Id:          claimData.Get("signing_channel.claim_id").String(),
			LbryUrl:     claimData.Get("signing_channel.canonical_url").String(),
			Description: template.HTML(claimData.Get("signing_channel.value.description").String()),
			Thumbnail:   utils.ToProxiedImageUrl(claimData.Get("signing_channel.value.thumbnail.url").String()),
		},
		Duration:     utils.FormatDuration(claimData.Get("value.video.duration").Int()),
		Title:        claimData.Get("value.title").String(),
		Description:  template.HTML(utils.ProcessText(claimData.Get("value.description").String(), true)),
		ThumbnailUrl: utils.ToProxiedImageUrl(claimData.Get("value.thumbnail.url").String()),
		License:      claimData.Get("value.license").String(),
		ValueType:    claimData.Get("value_type").String(),
		Repost:       claimData.Get("reposted_claim.canonical_url").String(),
		MediaType:    claimData.Get("value.source.media_type").String(),
		StreamType:   claimData.Get("value.stream_type").String(),
		HasFee:       claimData.Get("value.fee").Exists(),
	}

	timestamp := claimData.Get("value.release_time").Int()
	if timestamp == 0 {
		timestamp = claimData.Get("timestamp").Int()
	}
	claim.Time = time.Unix(timestamp, 0)
	claim.RelTime = humanize.Time(claim.Time)
	claim.Date = claim.Time.Format("January 2, 2006")

	claimData.Get("value.tags").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			claim.Tags = append(claim.Tags, value.String())
			return true
		},
	)

	url, err := utils.LbryTo(claim.LbryUrl)
	if err != nil {
		return Claim{}, err
	}
	claim.Url = url["http"]
	claim.RelUrl = url["rel"]
	claim.OdyseeUrl = url["odysee"]

	channelUrl, err := utils.LbryTo(claim.Channel.LbryUrl)
	if err != nil {
		return Claim{}, err
	}
	claim.Channel.Url = channelUrl["http"]
	claim.Channel.RelUrl = channelUrl["rel"]
	claim.Channel.OdyseeUrl = channelUrl["odysee"]

	return claim, nil
}

func (claim *Claim) GetViews() error {
	cacheData, found := claimCache.Get(claim.Id + "-views")
	if found {
		claim.Views = cacheData.(int64)
		return nil
	}

	data, err := utils.RequestJSON("https://api.odysee.com/file/view_count?auth_token="+viper.GetString("AUTH_TOKEN")+"&claim_id="+claim.Id, nil)
	if err != nil {
		return err
	}

	claim.Views = data.Get("data.0").Int()
	claimCache.Set(claim.Id+"-views", claim.Views, cache.DefaultExpiration)
	return nil
}

func (claim *Claim) GetRatings() error {
	cacheData, found := claimCache.Get(claim.Id + "-reactions")
	if found {
		ratings := cacheData.([]int64)
		claim.Likes = ratings[0]
		claim.Dislikes = ratings[1]
	}

	formData := url.Values{
		"claim_ids": []string{claim.Id},
	}
	body, err := utils.Request("https://api.odysee.com/reaction/list", 1000000, utils.Data{
		Bytes: strings.NewReader(formData.Encode()),
		Type:  "application/x-www-form-urlencoded",
	})
	if err != nil {
		return err
	}

	data := gjson.Parse(string(body))
	ratings := []int64{
		data.Get("data.others_reactions." + claim.Id + ".like").Int(),
		data.Get("data.others_reactions." + claim.Id + ".dislike").Int(),
	}
	claim.Likes = ratings[0]
	claim.Dislikes = ratings[1]

	claimCache.Set(claim.Id+"-reactions", ratings, cache.DefaultExpiration)
	return nil
}
