package main

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mguterl/phishin-discord-bot/phishin"
)

func formatDayOfWeek(d phishin.Date) string {
	return d.Format("Monday, January 2, 2006")
}

func formatDate(d phishin.Date) string {
	return fmt.Sprintf("%d/%d/%d", d.Month(), d.Day(), d.Year())
}

func formatDuration(t time.Duration) string {
	s := t.Seconds()
	hh := int(s / 3600)
	mm := int(math.Mod(s, 3600) / 60)
	ss := int(math.Mod(s, 60))

	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

func embedForSetlist(setlist *Setlist) (discordgo.MessageEmbed, error) {
	var d bytes.Buffer
	for _, set := range setlist.Sets {
		d.WriteString(fmt.Sprintf("🔥 __%s__ (%s)\n", set.Name, formatDuration(set.Duration)))

		for i, song := range set.Songs {
			d.WriteString(fmt.Sprintf("**%d.** %s", i+1, song))
			if i != len(set.Songs)-1 {
				d.WriteString(" ")
			}
		}
		d.WriteString("\n\n")
	}
	d.WriteString(fmt.Sprintf("%s\n", phishin.ShowURL(setlist.Date)))

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       fmt.Sprintf("Setlist for %s @ %s in %s", formatDate(setlist.Date), setlist.Venue, setlist.Location),
		Description: d.String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Total set duration: %s", formatDuration(setlist.Duration)),
		},
	}, nil
}

func embedForLastPlayedTracks(lastPlayed *phishin.LastPlayedTracksResponse) (discordgo.MessageEmbed, error) {
	var d bytes.Buffer
	last, rest := lastPlayed.Shows[len(lastPlayed.Shows)-1], lastPlayed.Shows[:len(lastPlayed.Shows)-1]

	d.WriteString(fmt.Sprintf("It was played at %s in %s\n", last.Venue.Name, last.Venue.Location))

	if len(rest) > 0 {
		d.WriteString("\nNext most recent plays 🌸:\n")
		for i := len(rest) - 1; i >= 0; i-- {
			show := rest[i]
			d.WriteString(fmt.Sprintf("🌵 %s at %s in %s\n", formatDayOfWeek(show.Date), show.Venue.Name, show.Venue.Location))
		}
	}
	d.WriteString(fmt.Sprintf("\n%s\n", lastPlayed.URL))

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       fmt.Sprintf("%s was last played on %s", lastPlayed.Title, formatDayOfWeek(last.Date)),
		Description: d.String(),
	}, nil
}

func embedForLongestTracks(longest *phishin.LongestTracksResponse) (discordgo.MessageEmbed, error) {
	var title string

	if len(longest.Tracks) == 1 {
		title = fmt.Sprintf("The longest version of %s", longest.Title)
	} else {
		title = fmt.Sprintf("The %d longest versions of %s", len(longest.Tracks), longest.Title)
	}

	var d bytes.Buffer
	for _, track := range longest.Tracks {
		d.WriteString(fmt.Sprintf("%s on %s 🐟🐟🐟 @ %s in %s\n\n", formatDuration(track.Duration), formatDayOfWeek(track.Show.Date), track.Show.Venue.Name, track.Show.Venue.Location))
	}

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       title,
		Description: d.String(),
	}, nil
}
