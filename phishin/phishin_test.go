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

func TestSongs(t *testing.T) {
	token := os.Getenv("PHISHIN_TOKEN")
	p := New(token)
	song, err := p.SongByTitle(context.Background(), "Scent of a Mule")
	require.NoError(t, err)
	assert.True(t, song.Success)
	assert.Equal(t, "Scent of a Mule", song.Data.Title)
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
}
