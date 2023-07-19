package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var channelCache = cache.New(30*time.Minute, 15*time.Minute)

type Channel struct {
	Name           string
	Title          string
	Id             string
	Followers      int64
	Url            string
	RelUrl         string
	OdyseeUrl      string
	LbryUrl        string
	CoverImg       string
	Description    template.HTML
	DescriptionTxt string
	Thumbnail      string
	ValueType      string
	UploadCount    int64
	Claims				 []Claim
}

func GetChannel(channel string) (Channel, error) {
	cacheData, found := channelCache.Get(channel)
	if found {
		return cacheData.(Channel), nil
	}

	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     []string{channel},
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=resolve", bytes.NewBuffer(resolveData))
	if err != nil {
		return Channel{}, err
	}
	data = data.Get("result." + strings.ReplaceAll(channel, ".", "\\."))

	channelData, err := ProcessChannel(data)
	if err != nil {
		return Channel{}, err
	}

	channelCache.Set(channel, channelData, cache.DefaultExpiration)
	return channelData, nil
}

func ProcessChannel(data gjson.Result) (Channel, error) {
	channel := Channel{
		Name:           data.Get("name").String(),
		Title:          data.Get("value.title").String(),
		Id:             data.Get("claim_id").String(),
		LbryUrl:				data.Get("canonical_url").String(),
		CoverImg:       utils.ToProxiedImageUrl(data.Get("value.cover.url").String()),
		Description:    template.HTML(utils.ProcessText(data.Get("value.description").String(), true)),
		DescriptionTxt: bluemonday.StrictPolicy().Sanitize(data.Get("value.description").String()),
		Thumbnail:      utils.ToProxiedImageUrl(data.Get("value.thumbnail.url").String()),
		ValueType:      data.Get("value_type").String(),
		UploadCount:    data.Get("meta.claims_in_channel").Int(),
	}

	url, err := utils.LbryTo(data.Get("canonical_url").String())
	if err != nil {
		return Channel{}, err
	}
	channel.Url = url["http"]
	channel.RelUrl = url["rel"]
	channel.OdyseeUrl = url["odysee"]

	return channel, nil
}

func (channel *Channel) GetFollowers() (int64, error) {
	cacheData, found := channelCache.Get(channel.Id + "-followers")
	if found {
		channel.Followers = cacheData.(int64)
		return channel.Followers, nil
	}

	data, err := utils.RequestJSON("https://api.odysee.com/subscription/sub_count?auth_token="+viper.GetString("AUTH_TOKEN")+"&claim_id="+channel.Id, nil)
	if err != nil {
		return 0, err
	}

	channel.Followers = data.Get("data.0").Int()
	channelCache.Set(channel.Id+"-followers", channel.Followers, cache.DefaultExpiration)
	return channel.Followers, err
}

func (channel Channel) GetClaims(page int) ([]Claim, error) {
	cacheData, found := channelCache.Get(channel.Id + "-claims-" + fmt.Sprint(page))
	if found {
		return cacheData.([]Claim), nil
	}

	channelDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":                20,
			"page":                     page,
			"no_totals":                true,
			"claim_type":               []string{"stream"},
			"order_by":                 []string{"release_time"},
			"fee_amount":               "<=0",
			"channel_ids":              []string{channel.Id},
			"release_time":             "<" + fmt.Sprint(time.Now().Unix()),
			"include_purchase_receipt": true,
		},
	}
	channelData, _ := json.Marshal(channelDataMap)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=claim_search", bytes.NewBuffer(channelData))
	if err != nil {
		return []Claim{}, nil
	}

	claims := make([]Claim, 0)
	claimIds := map[string]int64{}
	wg := sync.WaitGroup{}
	data.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			claimIds[value.Get("claim_id").String()] = key.Int()

			wg.Add(1)
			go func() {
				defer wg.Done()
				claim, _ := ProcessClaim(value)
				claim.GetViews()
				claims = append(claims, claim)
			}()

			return true
		},
	)
	wg.Wait()

	sort.Slice(claims, func(i, j int) bool {
		return claimIds[claims[i].Id] < claimIds[claims[j].Id]
	})

	channel.Claims = claims

	channelCache.Set(channel.Id+"-claims-"+fmt.Sprint(page), claims, cache.DefaultExpiration)
	return claims, nil
}
