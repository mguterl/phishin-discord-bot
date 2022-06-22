package phishin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseUrl = "http://phish.in/api/v1"

type Client struct {
	token string
	http  http.Client
}

func New(token string) Client {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	return Client{token: token, http: client}
}

type Venue struct {
	ID         int           `json:"id"`
	Slug       string        `json:"slug"`
	Name       string        `json:"name"`
	OtherNames []interface{} `json:"other_names"`
	Latitude   float64       `json:"latitude"`
	Longitude  float64       `json:"longitude"`
	ShowsCount int           `json:"shows_count"`
	Location   string        `json:"location"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type Track struct {
	ID                int           `json:"id"`
	ShowID            int           `json:"show_id"`
	ShowDate          string        `json:"show_date"`
	VenueName         string        `json:"venue_name"`
	VenueLocation     string        `json:"venue_location"`
	Title             string        `json:"title"`
	Position          int           `json:"position"`
	Duration          int           `json:"duration"`
	JamStartsAtSecond interface{}   `json:"jam_starts_at_second"`
	Set               string        `json:"set"`
	SetName           string        `json:"set_name"`
	LikesCount        int           `json:"likes_count"`
	Slug              string        `json:"slug"`
	Tags              []interface{} `json:"tags"`
	Mp3               string        `json:"mp3"`
	WaveformImage     string        `json:"waveform_image"`
	SongIds           []int         `json:"song_ids"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type Show struct {
	ID         int           `json:"id"`
	Date       string        `json:"date"`
	Duration   int           `json:"duration"`
	Incomplete bool          `json:"incomplete"`
	Sbd        bool          `json:"sbd"`
	Remastered bool          `json:"remastered"`
	Tags       []interface{} `json:"tags"`
	TourID     int           `json:"tour_id"`
	Venue      Venue         `json:"venue"`
	VenueName  string        `json:"venue_name"`
	TaperNotes string        `json:"taper_notes"`
	LikesCount int           `json:"likes_count"`
	Tracks     []Track       `json:"tracks"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type ShowOnDate struct {
	Success      bool `json:"success"`
	TotalEntries int  `json:"total_entries"`
	TotalPages   int  `json:"total_pages"`
	Page         int  `json:"page"`
	Data         Show `json:"data"`
}

func (c Client) ShowOnDate(ctx context.Context, t time.Time) (ShowOnDate, error) {
	date := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	url := fmt.Sprintf("%s/%s/%s", baseUrl, "show-on-date", date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ShowOnDate{}, fmt.Errorf("ShowOnDate: %v: %w", date, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.http.Do(req)
	if err != nil {
		return ShowOnDate{}, fmt.Errorf("ShowOnDate: %v: %w", date, err)
	}
	defer resp.Body.Close()

	var s ShowOnDate
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return s, fmt.Errorf("ShowOnDate: %v: %w", date, err)
	}

	return s, nil
}
