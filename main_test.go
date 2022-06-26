package main

import (
	"testing"
	"time"

	"github.com/mguterl/phishin-discord-bot/phishin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseDuration(t *testing.T, s string) time.Duration {
	duration, err := time.ParseDuration(s)
	require.NoError(t, err)
	return duration
}
func TestSetlist(t *testing.T) {
	var hourInMillis int = 3.6e+6 // 1 hour in millis

	show := phishin.Show{
		Date: phishin.DateFromString("2021-08-31"),
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
	assert.Equal(t, phishin.DateFromString("2021-08-31"), setlist.Date)
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

func TestParseCommand(t *testing.T) {
	command, args, err := parseCommand(".setlist 8/31/21")
	require.NoError(t, err)
	assert.Equal(t, SetlistCommand, command)
	d, ok := args.(time.Time)
	require.True(t, ok)
	assert.Equal(t, date(2021, 8, 31), d)
}
