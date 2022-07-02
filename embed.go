package main

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mguterl/phishin-discord-bot/phishin"
)

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
	d.WriteString(fmt.Sprintf("https://phish.in/%s\n", setlist.Date.Format("2006-01-02")))

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       fmt.Sprintf("Setlist for %s @ %s in %s", formatDate(setlist.Date.Time), setlist.Venue, setlist.Location),
		Description: d.String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Total set duration: %s", formatDuration(setlist.Duration)),
		},
	}, nil
}

func embedForLastPlayed(lastPlayed *phishin.LastPlayed) (discordgo.MessageEmbed, error) {
	var d bytes.Buffer
	last, rest := lastPlayed.Shows[len(lastPlayed.Shows)-1], lastPlayed.Shows[:len(lastPlayed.Shows)-1]

	d.WriteString(fmt.Sprintf("It was played at %s in %s\n", last.Venue.Name, last.Venue.Location))

	if len(rest) > 0 {
		d.WriteString("\nNext most recent plays 🌸:\n")
		for i := len(rest) - 1; i >= 0; i-- {
			show := rest[i]
			d.WriteString(fmt.Sprintf("🌵 %s at %s in %s\n", show.Date.Format("Monday, January 2, 2006"), show.Venue.Name, show.Venue.Location))
		}
	}
	d.WriteString(fmt.Sprintf("\n%s\n", lastPlayed.URL))

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       fmt.Sprintf("%s was last played on %s", lastPlayed.Title, last.Date.Format("Monday, January 2, 2006")),
		Description: d.String(),
	}, nil
}

func embedForLongest(longest *phishin.Longest) (discordgo.MessageEmbed, error) {
	var title string

	if len(longest.Tracks) == 1 {
		title = fmt.Sprintf("The longest version of %s", longest.Title)
	} else {
		title = fmt.Sprintf("The %d longest versions of %s", len(longest.Tracks), longest.Title)
	}

	var d bytes.Buffer
	for _, track := range longest.Tracks {
		d.WriteString(fmt.Sprintf("%s on %s 🐟🐟🐟 @ %s in %s\n\n", formatDuration(track.Duration), track.Show.Date.Format("Monday, January 2, 2006"), track.Show.Venue.Name, track.Show.Venue.Location))
	}

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       title,
		Description: d.String(),
	}, nil
}
