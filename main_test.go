package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestParseSetlistCommand(t *testing.T) {
	command, args, err := parseCommand(".setlist 8/31/21")
	require.NoError(t, err)
	assert.Equal(t, SetlistCommand, command)
	d, ok := args.(time.Time)
	require.True(t, ok)
	assert.Equal(t, date(2021, 8, 31), d)
}

func TestParseSetlistCommandError(t *testing.T) {
	command, _, err := parseCommand(".setlist notadate")
	assert.Error(t, err)
	assert.Equal(t, SetlistCommand, command)
}

func TestParseRandomCommand(t *testing.T) {
	command, _, err := parseCommand(".random")
	require.NoError(t, err)
	assert.Equal(t, RandomCommand, command)
}

func TestParseLastPlayedCommand(t *testing.T) {
	command, args, err := parseCommand(".lastplayed song title")
	require.NoError(t, err)
	assert.Equal(t, LastPlayedCommand, command)
	s, ok := args.(string)
	require.True(t, ok)
	assert.Equal(t, "song title", s)
}
