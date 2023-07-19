package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/tidwall/gjson"
)

func Search(query string, page int, claimType string, nsfw bool, size int) ([]interface{}, error) {
	from := 0
	if page > 1 {
		from = page * size
	}

	query = strings.ReplaceAll(query, " ", "+")
	url := "https://lighthouse.odysee.tv/search?s=" + query + "&size=" + fmt.Sprint(size) + "&free_only=true&from=" + fmt.Sprint(from) + "&nsfw=" + strconv.FormatBool(nsfw) + "&claimType=" + claimType

	data, err := utils.RequestJSON(url, nil)
	if err != nil {
		return nil, err
	}

	results, err := ProcessResults(data, "")
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetRelated(query string, nsfw bool, relatedTo string) ([]interface{}, error) {
	query = strings.ReplaceAll(query, " ", "+")
	url := "https://recsys.odysee.tv/search?s=" + query + "&size=20" + "&from=0" + "&related_to=" + relatedTo + "&free_only=true" + "&nsfw=" + strconv.FormatBool(nsfw)

	data, err := utils.RequestJSON(url, nil)
	if err != nil {
		return nil, err
	}

	results, err := ProcessResults(data, relatedTo)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ProcessResults(data gjson.Result, relatedTo string) ([]interface{}, error) {
	results := make([]interface{}, 0)
	wg := sync.WaitGroup{}

	urls := []string{}
	data.ForEach(func(key gjson.Result, value gjson.Result) bool {
		urls = append(urls, "lbry://" + value.Get("name").String() + "#" + value.Get("claimId").String())
		return true
	})

	claims, err := GetClaims(urls, true, true)
	if err != nil {
		return nil, err
	}
	for _, claim := range claims {
		id := reflect.ValueOf(claim).FieldByName("Id").String()
		if err == nil && id != relatedTo {
			results = append(results, claim)
		}
	}

	wg.Wait()
	return results, nil
}
