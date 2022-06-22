package main

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mguterl/phishin-discord-bot/phishin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseDuration(t *testing.T, s string) time.Duration {
	duration, err := time.ParseDuration(s)
	require.NoError(t, err)
	return duration
}

func TestSetlistResponse(t *testing.T) {
	setlist := Setlist{
		Date:     parseDate("2021-08-31"),
		Venue:    "Shoreline Amphitheatre",
		Location: "Mountain View, CA",
		Duration: parseDuration(t, "3h6m45s"),
		Sets: []*Set{
			{
				Name:     "Set 1",
				Duration: parseDuration(t, "1h26m48s"),
				Songs:    []string{"Glide", "Colonel Forbin's Ascent"},
			},
			{
				Name:     "Set 2",
				Duration: parseDuration(t, "1h29m24s"),
				Songs:    []string{"Soul Planet"},
			},
			{
				Name:     "Encore",
				Duration: parseDuration(t, "10m32s"),
				Songs:    []string{"Fee", "Wilson"},
			},
		},
	}
	embed, err := embedForSetlist(setlist)
	require.NoError(t, err)
	assert.Equal(t, discordgo.MessageEmbed{
		Color: green,
		Title: "Setlist for 8/31/2021 @ Shoreline Amphitheatre in Mountain View, CA",
		Description: `🔥 __Set 1__ (01:26:48)
**1.** Glide **2.** Colonel Forbin's Ascent

🔥 __Set 2__ (01:29:24)
**1.** Soul Planet

🔥 __Encore__ (00:10:32)
**1.** Fee **2.** Wilson

https://phish.in/2021-08-31
`,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Total set duration: 03:06:45",
		},
	}, embed)
}

func TestSetlist(t *testing.T) {
	var hourInMillis int = 3.6e+6 // 1 hour in millis

	show := phishin.Show{
		Date: "2021-08-31",
		Venue: phishin.Venue{
			Name:     "Shoreline Amphitheatre",
			Location: "Mountain View, CA",
		},
		Duration: hourInMillis,
		Tracks: []phishin.Track{
			{
				SetName:  "Set 1",
				Title:    "Glide",
				Duration: 120000, // 2 minutes
			},
			{
				SetName:  "Set 1",
				Title:    "Colonel Forbin's Ascent",
				Duration: 180000, // 3 minutes
			},
			{
				SetName:  "Set 2",
				Title:    "Soul Planet",
				Duration: 240000, // 4 minutes
			},
			{
				SetName:  "Encore",
				Title:    "Fee",
				Duration: 60000, // 1 minute
			},
			{
				SetName:  "Encore",
				Title:    "Wilson",
				Duration: 120000, // 2 minutes
			},
		},
	}
	setlist := setlistForShow(show)
	assert.Equal(t, time.Date(2021, 8, 31, 0, 0, 0, 0, time.UTC), setlist.Date)
	assert.Equal(t, "Shoreline Amphitheatre", setlist.Venue)
	assert.Equal(t, "Mountain View, CA", setlist.Location)
	assert.Equal(t, 1*time.Hour, setlist.Duration)
	assert.Equal(t, []*Set{
		{
			Name:     "Set 1",
			Duration: 5 * time.Minute,
			Songs:    []string{"Glide", "Colonel Forbin's Ascent"},
		},
		{
			Name:     "Set 2",
			Duration: 4 * time.Minute,
			Songs:    []string{"Soul Planet"},
		},
		{
			Name:     "Encore",
			Duration: 3 * time.Minute,
			Songs:    []string{"Fee", "Wilson"},
		},
	}, setlist.Sets)
}

// Time returns a pointer value for the time.Time value passed in.
func Time(v time.Time) *time.Time {
	return &v
}

func tParseCommand(t *testing.T, s string) *time.Time {
	d, err := parseCommand(s)
	require.NoError(t, err)
	return d
}

func TestParseCommand(t *testing.T) {
	assert.Equal(t, Time(date(2021, 8, 31)), tParseCommand(t, ".setlist 2021-08-31"))
	assert.Equal(t, Time(date(2021, 8, 31)), tParseCommand(t, ".setlist 8/31/2021"))
	assert.Equal(t, Time(date(2021, 8, 31)), tParseCommand(t, ".setlist 8/31/21"))

	_, err := parseCommand(".setlist 42")
	require.Error(t, err)
}
