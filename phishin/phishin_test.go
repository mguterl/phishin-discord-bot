package phishin

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowOnDate(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	date := time.Date(2021, 8, 31, 0, 0, 0, 0, time.UTC)
	show, err := p.ShowOnDate(context.Background(), date)
	require.NoError(t, err)
	assert.Equal(t, "Shoreline Amphitheatre", show.Data.Venue.Name)
	assert.Len(t, show.Data.Tracks, 20)
}

func TestRandomShow(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	show, err := p.RandomShow(context.Background())
	require.NoError(t, err)
	assert.Greater(t, len(show.Data.Venue.Name), 0)
	assert.Greater(t, len(show.Data.Tracks), 5)
}

func TestSongByTitle(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	song, err := p.SongByTitle(context.Background(), "Scent of a Mule")
	require.NoError(t, err)
	assert.True(t, song.Success)
	assert.Equal(t, "Scent of a Mule", song.Data.Title)
}

func TestSongByTitleDoesNotExist(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	_, err := p.SongByTitle(context.Background(), "Does not exist")
	require.Error(t, err)
}

func TestLastPlayed(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	last, err := p.LastPlayedTracks(context.Background(), "Black-Eyed Katy", 4)
	require.NoError(t, err)
	assert.Equal(t, "Black-Eyed Katy", last.Title)
	require.Len(t, last.Shows, 4)
	mostRecent := last.Shows[len(last.Shows)-1]
	assert.Equal(t, DateFromString("1997-12-30"), mostRecent.Date)
	assert.Equal(t, "https://phish.in/black-eyed-katy", last.URL)
	assert.Equal(t, 7, last.PlayCount)
}

func TestLastPlayedMultipleSameShow(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	last, err := p.LastPlayedTracks(context.Background(), "Moby Dick", 3)
	require.NoError(t, err)
	assert.Equal(t, "Moby Dick", last.Title)
	require.Len(t, last.Shows, 3)
	first := last.Shows[0]
	assert.Equal(t, DateFromString("1993-02-19"), first.Date)
	assert.Equal(t, "https://phish.in/moby-dick", last.URL)
	assert.Equal(t, 3, last.PlayCount)
}

func TestLastPlayedMissingSong(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	_, err := p.LastPlayedTracks(context.Background(), "Does Not Exist", 3)
	require.Error(t, err)
}

func TestLongest(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	longest, err := p.LongestTracks(context.Background(), "Access Me", 3)
	require.NoError(t, err)
	assert.Equal(t, "Access Me", longest.Title)
	assert.Len(t, longest.Tracks, 3)
	assert.Equal(t, 6, longest.PlayCount)
}

func TestSlugify(t *testing.T) {
	assert.Equal(t, "black-eyed-katy", slugify("Black-Eyed Katy"))
	assert.Equal(t, "colonel-forbin-s-ascent", slugify("Colonel Forbin's Ascent"))
	assert.Equal(t, "ac-dc-bag", slugify("AC/DC Bag"))
	assert.Equal(t, "bill-bailey-won-t-you-please-come-home", slugify("Bill Bailey, Won't You Please Come Home?"))
	assert.Equal(t, "mean-mr-mustard", slugify("Mean Mr. Mustard"))
}
