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

func TestSetlistEmbed(t *testing.T) {
	setlist := &Setlist{
		Date:     phishin.DateFromString("2021-08-31"),
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

func TestLastPlayedEmbed(t *testing.T) {
	lastPlayed := &phishin.LastPlayedTracksResponse{
		Title:     "Black-Eyed Katy",
		URL:       "https://phish.in/black-eyed-katy",
		PlayCount: 7,
		Shows: []phishin.Show{
			{
				Date: phishin.DateFromString("1997-11-28"),
				Venue: phishin.Venue{
					Name:     "The Centrum",
					Location: "Worcester, MA",
				},
			},
			{
				Date: phishin.DateFromString("1997-12-05"),
				Venue: phishin.Venue{
					Name:     "CSU Convocation Center",
					Location: "Cleveland, OH",
				},
			},
			{
				Date: phishin.DateFromString("1997-12-30"),
				Venue: phishin.Venue{
					Name:     "Madison Square Garden",
					Location: "New York, NY",
				},
			},
		},
	}
	embed, err := embedForLastPlayedTracks(lastPlayed)
	require.NoError(t, err)
	assert.Equal(t, discordgo.MessageEmbed{
		Color: green,
		Title: "Black-Eyed Katy was last played on Tuesday, December 30, 1997",
		Description: `It was played at Madison Square Garden in New York, NY

Next most recent plays 🌸:
🌵 Friday, December 5, 1997 at CSU Convocation Center in Cleveland, OH
🌵 Friday, November 28, 1997 at The Centrum in Worcester, MA

https://phish.in/black-eyed-katy
`,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Total play count: 7",
		},
	}, embed)
}

func TestLastPlayedEmbedWithOnePlay(t *testing.T) {
	lastPlayed := &phishin.LastPlayedTracksResponse{
		Title: "And So To Bed",
		URL:   "https://phish.in/and-so-to-bed",
		Shows: []phishin.Show{
			{
				Date: phishin.DateFromString("2021-10-15"),
				Venue: phishin.Venue{
					Name:     "Golden 1 Center",
					Location: "Sacramento, CA",
				},
			},
		},
	}
	embed, err := embedForLastPlayedTracks(lastPlayed)
	require.NoError(t, err)
	assert.Equal(t, discordgo.MessageEmbed{
		Color: green,
		Title: "And So To Bed was last played on Friday, October 15, 2021",
		Description: `It was played at Golden 1 Center in Sacramento, CA

https://phish.in/and-so-to-bed
`, Footer: &discordgo.MessageEmbedFooter{
			Text: "Total play count: 0",
		},
	}, embed)
}
