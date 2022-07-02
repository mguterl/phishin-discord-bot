package phishin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

const baseApiUrl = "http://phish.in/api/v1"
const baseURL = "https://phish.in"

type Client struct {
	token string
	http  *http.Client
}

func New(token string) *Client {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	return &Client{token: token, http: &client}
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
	ShowDate          Date          `json:"show_date"`
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
	Date       Date          `json:"date"`
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

type ShowResponse struct {
	Success      bool `json:"success"`
	TotalEntries int  `json:"total_entries"`
	TotalPages   int  `json:"total_pages"`
	Page         int  `json:"page"`
	Data         Show `json:"data"`
}

type SongResponse struct {
	Success      bool `json:"success"`
	TotalEntries int  `json:"total_entries"`
	TotalPages   int  `json:"total_pages"`
	Page         int  `json:"page"`
	Data         Song `json:"data"`
}

type Song struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Alias       interface{} `json:"alias"`
	TracksCount int         `json:"tracks_count"`
	Slug        string      `json:"slug"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Tracks      []SongTrack `json:"tracks"`
}

func (s Song) URL() string {
	return fmt.Sprintf("%s/%s", baseURL, s.Slug)
}

type SongTrack struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Duration   int    `json:"duration"`
	ShowID     int    `json:"show_id"`
	ShowDate   Date   `json:"show_date"`
	Set        string `json:"set"`
	Position   int    `json:"position"`
	LikesCount int    `json:"likes_count"`
	Slug       string `json:"slug"`
	Mp3        string `json:"mp3"`
}

type Date struct {
	time.Time
}

type LastPlayedTracksResponse struct {
	Title string
	URL   string
	Shows []Show
}

type LongestTracksResponse struct {
	Title  string
	Tracks []LongestTrack
}

type LongestTrack struct {
	Duration time.Duration
	Show     Show
}

func DateFromTime(t time.Time) Date {
	return Date{Time: t}
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

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

var replacer = strings.NewReplacer(
	" ", "-",
	"'", "-",
	"/", "-",
	",", "",
	".", "",
	"?", "",
)

func slugify(s string) string {
	return strings.ToLower(replacer.Replace(s))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func duration(millis int) time.Duration {
	return time.Duration(millis * int(time.Millisecond))
}

func (c *Client) get(ctx context.Context, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil
	}

	if v != nil {
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

func ShowURL(d Date) string {
	return fmt.Sprintf("%s/%s", baseURL, d.Format("2006-01-02"))
}

func (c *Client) ShowOnDate(ctx context.Context, t time.Time) (*ShowResponse, error) {
	date := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	url := fmt.Sprintf("%s/%s/%s", baseApiUrl, "show-on-date", date)
	s := &ShowResponse{}
	return s, c.get(ctx, url, s)
}

func (c *Client) RandomShow(ctx context.Context) (*ShowResponse, error) {
	url := fmt.Sprintf("%s/%s", baseApiUrl, "random-show")
	s := &ShowResponse{}
	return s, c.get(ctx, url, s)
}

func (c *Client) SongByTitle(ctx context.Context, title string) (*SongResponse, error) {
	slug := slugify(title)
	url := fmt.Sprintf("%s/%s/%s", baseApiUrl, "songs", slug)
	s := &SongResponse{}
	err := c.get(ctx, url, s)
	if err != nil {
		return nil, err
	}
	if !s.Success { // TODO: Move this into get helper
		return nil, fmt.Errorf("SongByTitle not found: %v", title)
	}
	return s, nil
}

func (c *Client) LastPlayedTracks(ctx context.Context, title string, count int) (*LastPlayedTracksResponse, error) {
	song, err := c.SongByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	lastPlayed := &LastPlayedTracksResponse{
		Title: song.Data.Title,
		URL:   song.Data.URL(),
	}

	// A song may appear multiple times for a single show, but we only want to
	// display it once per show.
	var exists struct{}
	seen := make(map[int]struct{})
	tracks := make([]SongTrack, 0, count)

	// Iterate backwards through the list.
	for i := len(song.Data.Tracks) - 1; i >= 0 && len(tracks) < count; i-- {
		track := song.Data.Tracks[i]
		if _, ok := seen[track.ShowID]; !ok {
			tracks = append(tracks, track)
			seen[track.ShowID] = exists
		}
	}

	// Keep shows in date ascending order to match the original API response.
	for i := len(tracks) - 1; i >= 0; i-- {
		track := tracks[i]
		show, err := c.ShowOnDate(ctx, track.ShowDate.Time)
		if err != nil {
			return nil, err
		}
		lastPlayed.Shows = append(lastPlayed.Shows, show.Data)
	}

	return lastPlayed, nil
}

func (c *Client) LongestTracks(ctx context.Context, title string, count int) (*LongestTracksResponse, error) {
	song, err := c.SongByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("longest tracks: %v: %w", title, err)
	}

	longest := &LongestTracksResponse{
		Title: song.Data.Title,
	}
	tracks := song.Data.Tracks
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Duration > tracks[j].Duration
	})

	for i := 0; i < min(len(tracks), count); i++ {
		track := tracks[i]
		longestTrack := LongestTrack{
			Duration: duration(track.Duration),
		}

		show, err := c.ShowOnDate(ctx, track.ShowDate.Time)
		if err != nil {
			return longest, err
		}

		longestTrack.Show = show.Data
		longest.Tracks = append(longest.Tracks, longestTrack)
	}

	return longest, nil
}
