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
	assert.Equal(t, "Shoreline Amphitheatre", show.Data.Venue.Name)
	assert.Len(t, show.Data.Tracks, 20)
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
	last, err := p.LastPlayed(context.Background(), "Black-Eyed Katy", 4)
	require.NoError(t, err)
	assert.Equal(t, "Black-Eyed Katy", last.Title)
	require.Len(t, last.Shows, 4)
	mostRecent := last.Shows[len(last.Shows)-1]
	assert.Equal(t, DateFromString("1997-12-30"), mostRecent.Date)
	assert.Equal(t, "https://phish.in/black-eyed-katy", last.URL)
}

func TestLastPlayedMultipleSameShow(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	last, err := p.LastPlayed(context.Background(), "Moby Dick", 3)
	require.NoError(t, err)
	assert.Equal(t, "Moby Dick", last.Title)
	require.Len(t, last.Shows, 3)
	first := last.Shows[0]
	assert.Equal(t, DateFromString("1993-02-19"), first.Date)
	assert.Equal(t, "https://phish.in/moby-dick", last.URL)
}

func TestLastPlayedMissingSong(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	_, err := p.LastPlayed(context.Background(), "Does Not Exist", 3)
	require.Error(t, err)
}
