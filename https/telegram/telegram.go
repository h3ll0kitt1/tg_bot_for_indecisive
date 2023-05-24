package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	pathAPI = "https://api.telegram.org/bot"
)

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

type Counter struct {
	offset int
}

func (c *Counter) NextUpdate(limit int, token string) (*UpdatesResponse, error) {

	u, err := getUpdatesFromBot(c.offset, limit, token)
	if err != nil {
		return nil, fmt.Errorf("GET next update: %w", err)
	}

	if len(u.Result) != 0 {
		c.offset = u.Result[0].UpdateId + 1
		return u, nil
	}
	return nil, nil
}

func SendMessageToChat(u UpdatesResponse, token string) (string, error) {

	baseURL := pathAPI + token + "/sendMessage"

	resp, err := http.PostForm(baseURL, url.Values{
		"chat_id": {strconv.Itoa(u.Result[0].Message.Chat.Id)},
		"text":    {u.Result[0].Message.Text},
	})
	if err != nil {
		return "", fmt.Errorf("POST request to chat: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected Status Code: %w", err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Parse response: %w", err)
	}

	return string(result), nil
}

func getUpdatesFromBot(offset int, limit int, token string) (*UpdatesResponse, error) {

	baseURL, err := url.Parse(pathAPI + token + "/getUpdates")
	if err != nil {
		return nil, fmt.Errorf("URL parse: %w", err)
	}

	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("GET updates: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected Status Code: %w", err)
	}

	var u UpdatesResponse

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("Decode json to model: %w", err)
	}
	return &u, nil
}
