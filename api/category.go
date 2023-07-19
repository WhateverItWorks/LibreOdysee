package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/utils"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var fpCache = cache.New(30*time.Minute, 30*time.Minute)
var categoryCache = cache.New(72*time.Hour, 36*time.Hour)

type Category struct {
	Name          string
	Path					string
	ChannelIds    []string
	NotChannelIds []string
	ChannelLimit	int64
	DaysOfContent int64
	PageSize      int64
}

func GetCategoryData() (map[string]Category, error) {
	cacheData, found := categoryCache.Get("en")
	if found {
		return cacheData.(map[string]Category), nil
	}

	data, err := utils.RequestJSON("https://odysee.com/$/api/content/v2/get", nil)
	if err != nil {
		return map[string]Category{}, nil
	}

	categories := map[string]Category{}
	data.Get("data.en.categories").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			categories[value.Get("name").String()] = Category{
				Name:          value.Get("label").String(),
				Path: 				 value.Get("name").String(),
				NotChannelIds: strings.Split(strings.TrimSuffix(strings.TrimPrefix(value.Get("excludedChannelIds").Raw, `["`), `"]`), `","`),
				ChannelIds:    strings.Split(strings.TrimSuffix(strings.TrimPrefix(value.Get("channelIds").Raw, `["`), `"]`), `","`),
				ChannelLimit:  value.Get("channelLimit").Int(),
				DaysOfContent: value.Get("daysOfContent").Int(),
				PageSize:      value.Get("pageSize").Int(),
			}
			return true
		},
	)
	categoryCache.Set("en", categories, cache.DefaultExpiration)
	return categories, nil
}

func GetOrderedCategoriesArray() ([]Category, error) {
	categories, err := GetCategoryData()
	if err != nil {
		return []Category{}, err
	}

	newCategories := []Category{}
	newCategories = append(newCategories, categories["featured"])
	newCategories = append(newCategories, categories["popculture"])
	newCategories = append(newCategories, categories["artists"])
	newCategories = append(newCategories, categories["education"])
	newCategories = append(newCategories, categories["lifestyle"])
	newCategories = append(newCategories, categories["spooky"])
	newCategories = append(newCategories, categories["gaming"])
	newCategories = append(newCategories, categories["tech"])
	newCategories = append(newCategories, categories["comedy"])
	newCategories = append(newCategories, categories["music"])
	newCategories = append(newCategories, categories["sports"])
	newCategories = append(newCategories, categories["universe"])
	newCategories = append(newCategories, categories["finance"])
	newCategories = append(newCategories, categories["news"])
	newCategories = append(newCategories, categories["rabbithole"])

	return newCategories, nil
}

func (category Category) GetCategoryClaims(page int, nsfw bool) ([]interface{}, error) {
	cacheData, found := fpCache.Get(category.Name)
	if found {
		return cacheData.([]interface{}), nil
	}

	nsfwTags := []string{"porn", "porno", "nsfw", "mature", "xxx", "sex", "creampie", "blowjob", "handjob", "vagina", "boobs", "big boobs", "big dick", "pussy", "cumshot", "anal", "hard fucking", "ass", "fuck", "hentai"}
	if nsfw {
		nsfwTags = []string{}
	}

	claimSearchData := map[string]interface{}{}
	if category.Name != "Rabbit Hole" {
		claimSearchData = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "claim_search",
			"params": map[string]interface{}{
				"any_languages":						[]string{"en", "none"},
				"channel_ids":              category.ChannelIds,
				"claim_type":               []string{"stream", "channel"},
				"duration":									">=60",
				"has_source":               true,
				"limit_claims_per_channel": category.ChannelLimit,
				"no_totals":                true,
				"not_channel_ids":          category.NotChannelIds,
				"not_tags":                 nsfwTags,
				"order_by":                 []string{"effective_amount"},
				"page":                     page,
				"page_size":                category.PageSize,
				"fee_amount":               "<=0",
				"release_time":             ">" + fmt.Sprint(time.Now().Unix()-(category.DaysOfContent*24*60*60)),
				"remove_duplicates":        true,
			},
		}
	} else {
		claimSearchData = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "claim_search",
			"params": map[string]interface{}{
				"any_languages":						[]string{"en", "none"},
				"claim_type":               []string{"stream", "channel"},
				"has_source":               true,
				"limit_claims_per_channel": 3,
				"no_totals":                true,
				"not_channel_ids":          category.NotChannelIds,
				"not_tags":                 nsfwTags,
				"order_by":                 []string{"trending_group", "trending_mixed"},
				"page":                     page,
				"page_size":                36,
				"release_time":             ">" + fmt.Sprint(time.Now().Unix()-518400),
				"remove_duplicates":        true,
				"fee_amount":               "<=0",
			},
		}
	}
	claimSearchReqData, _ := json.Marshal(claimSearchData)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=claim_search", bytes.NewBuffer(claimSearchReqData))
	if err != nil {
		return []interface{}{}, err
	}

	claims := make([]interface{}, 0)
	wg := sync.WaitGroup{}
	data.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)
			go func() {
				defer wg.Done()

				claim, err := ProcessClaim(value)
				if err != nil {
					channel, _ := ProcessChannel(value)
					channel.GetFollowers()
					claims = append(claims, channel)
				}
				claim.GetViews()
				claims = append(claims, claim)
			}()

			return true
		},
	)
	wg.Wait()

	fpCache.Set(category.Name, claims, cache.DefaultExpiration)
	return claims, nil
}
