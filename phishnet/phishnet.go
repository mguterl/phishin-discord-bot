package phishnet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func New(apiKey string) *Client {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	return &Client{apiKey: apiKey, http: &client}
}

type ShowsResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
	Shows        []Show `json:"data"`
}

type Show struct {
	ID        string `json:"showid"`
	Date      Date   `json:"showdate"`
	Permalink string `json:"permalink"`
	Venue     string `json:"venue"`
	City      string `json:"city"`
	State     string `json:"state"`
}

func (c *Client) Shows(ctx context.Context) (*ShowsResponse, error) {
	url := fmt.Sprintf("https://api.phish.net/v5/shows/artist/phish.json?order_by=showdate&apikey=%s", c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	s := &ShowsResponse{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func DateFromString(s string) Date {
	d := Date{}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		d.Time = time.Time{}
	}
	d.Time = t
	return d
}

func NextShow(shows []Show, now time.Time) *Show {
	for _, s := range shows {
		if s.Date.After(now) {
			return &s
		}
	}
	return nil
}
