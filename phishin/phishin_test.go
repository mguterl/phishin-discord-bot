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
