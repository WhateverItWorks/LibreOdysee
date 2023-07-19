package api

import (
	"log"
	"net/url"
	"strings"

	"codeberg.org/librarian/librarian/utils"
	"github.com/tidwall/gjson"
)

func NewUser() string {
	formData := url.Values{
		"auth_token": []string{},
		"language": []string{"en"},
		"app_id": []string{"odyseecom692EAWhtoqDuAfQ6KHMXxFxt8tkhmt7sfprEMHWKjy5hf6PwZcHDV542V"},
	}
	body, err := utils.Request("https://api.odysee.com/user/new", 1000000, utils.Data{
		Bytes: strings.NewReader(formData.Encode()),
		Type: "application/x-www-form-urlencoded",
	})
	if err != nil {
		log.Fatal(err)
	}

	return gjson.Get(string(body), "data.auth_token").String()
}
