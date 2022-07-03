package phishnet

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShows(t *testing.T) {
	p := New(os.Getenv("PHISHNET_TOKEN"))
	resp, err := p.Shows(context.Background())
	require.NoError(t, err)
	require.False(t, resp.Error) // TODO: this should return error
	assert.Greater(t, len(resp.Shows), 42)
}

func TestNextShow(t *testing.T) {
	shows := []Show{
		{
			ID:   "1",
			Date: DateFromString("2022-06-05"),
		},
		{
			ID:   "2",
			Date: DateFromString("2022-07-14"),
		},
		{
			ID:   "3",
			Date: DateFromString("2022-07-15"),
		},
	}

	now := time.Date(2022, time.July, 3, 0, 0, 0, 0, time.UTC)
	next := NextShow(shows, now)
	assert.Equal(t, "2", next.ID)

	next = NextShow([]Show{}, now)
	assert.Nil(t, next)

	now = time.Date(2023, time.July, 3, 0, 0, 0, 0, time.UTC)
	next = NextShow(shows, now)
	assert.Nil(t, next)
}
