package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/WhateverItWorks/LibreOdysee/utils"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var commentCache = cache.New(30*time.Minute, 15*time.Minute)

type Comments struct {
	Comments []Comment
	Items    int64
	Pages    int64
}

type Comment struct {
	Channel   Channel
	Comment   template.HTML
	CommentId string
	ParentId  string
	Pinned    bool
	Time      string
	RelTime   string
	Replies   int64
	Likes     int64
	Dislikes  int64
}

func CommentsHandler(c *fiber.Ctx) error {
	claimId := c.Query("claim_id")
	channelId := c.Query("channel_id")
	channelName := c.Query("channel_name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	if claimId == "" || channelId == "" || channelName == "" || page == "" || pageSize == "" {
		_, err := c.Status(400).WriteString("missing query param. claim_id, channel_id, channel_name, page, page_size required")
		return err
	}

	sortBy := 3
	switch c.Query("sort_by") {
	case "controversial":
		sortBy = 2
	case "new":
		sortBy = 0
	}

	newPage, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	newPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		return err
	}

	claim := Claim{
		Id: claimId,
		Channel: Channel{
			Id:   channelId,
			Name: channelName,
		},
	}
	comments, err := claim.GetComments(c.Query("parent_id"), sortBy, newPageSize, newPage)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	return c.JSON(comments)
}

func (claim Claim) GetComments(parentId string, sortBy int, pageSize int, page int) (Comments, error) {
	cacheData, found := commentCache.Get(claim.Id + parentId + fmt.Sprint(sortBy) + fmt.Sprint(page) + fmt.Sprint(pageSize))
	if found {
		return cacheData.(Comments), nil
	}

	reqDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment.List",
		"params": map[string]interface{}{
			"page":         page,
			"claim_id":     claim.Id,
			"page_size":    pageSize,
			"sort_by":      sortBy,
			"top_level":    true,
			"channel_id":   claim.Channel.Id,
			"channel_name": claim.Channel.Name,
		},
	}
	if parentId != "" {
		reqDataMap["params"].(map[string]interface{})["parent_id"] = parentId
		reqDataMap["params"].(map[string]interface{})["top_level"] = false
	}

	reqData, err := json.Marshal(reqDataMap)
	if err != nil {
		return Comments{}, err
	}

	data, err := utils.RequestJSON("https://comments.odysee.tv/api/v2?m=comment.List", bytes.NewBuffer(reqData))
	if err != nil {
		return Comments{}, err
	}

	commentIds := []string{}
	data.Get("result.items.#.comment_id").ForEach(
		func(key, value gjson.Result) bool {
			commentIds = append(commentIds, value.String())
			return true
		},
	)
	likesDislikes := GetCommentLikeDislikes(commentIds)

	channelUrls := strings.Split(strings.Trim(strings.ReplaceAll(data.Get("result.items.#.channel_url").Raw, `"`, ""), "[]"), ",")
	channelsInt, err := GetClaims(channelUrls, false, false)
	if err != nil {
		return Comments{}, err
	}
	channels := map[string]Channel{}
	for _, channelInt := range channelsInt {
		//lint:ignore S1034 
		switch channelInt.(type) {
		case Channel:
			channel := channelInt.(Channel)
			channels[channel.Id] = channel
		}
	}

	comments := []Comment{}
	data.Get("result.items").ForEach(
		func(key, value gjson.Result) bool {
			comment := Comment{
				Comment:   template.HTML(utils.ProcessText(value.Get("comment").String(), false)),
				CommentId: value.Get("comment_id").String(),
				ParentId:  value.Get("parent_id").String(),
				Replies:   value.Get("replies").Int(),
				Pinned:    value.Get("is_pinned").Bool(),
			}

			timestamp := time.Unix(value.Get("timestamp").Int(), 0)
			comment.Time = timestamp.UTC().Format("January 2, 2006 15:04")
			comment.RelTime = humanize.Time(timestamp)
			if comment.RelTime == "a long while ago" {
				comment.RelTime = comment.Time
			}

			comment.Likes = likesDislikes[comment.CommentId][0]
			comment.Dislikes = likesDislikes[comment.CommentId][1]

			comment.Channel = channels[strings.Split(value.Get("channel_url").String(), "#")[1]]

			comments = append(comments, comment)

			return true
		},
	)

	returnData := Comments{
		Comments: comments,
		Pages:    data.Get("result.total_pages").Int(),
		Items:    data.Get("result.total_items").Int(),
	}

	commentCache.Set(claim.Id + parentId + fmt.Sprint(sortBy) + fmt.Sprint(page) + fmt.Sprint(pageSize), returnData, cache.DefaultExpiration)
	return returnData, nil
}

func GetCommentLikeDislikes(commentIds []string) map[string][]int64 {
	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "reaction.List",
		"params": map[string]interface{}{
			"comment_ids": strings.Join(commentIds, ","),
		},
	}
	commentsData, _ := json.Marshal(commentsDataMap)

	data, err := utils.RequestJSON("https://comments.odysee.tv/api/v2?m=reaction.List", bytes.NewBuffer(commentsData))
	if err != nil {
		fmt.Println(err)
	}

	likesDislikes := make(map[string][]int64)
	data.Get("result.others_reactions").ForEach(
		func(key, value gjson.Result) bool {
			likesDislikes[key.String()] = []int64{
				value.Get("like").Int(),
				value.Get("dislike").Int(),
			}
			return true
		},
	)

	return likesDislikes
}
