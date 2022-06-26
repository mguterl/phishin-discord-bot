package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mguterl/phishin-discord-bot/phishin"
)

func main() {
	godotenv.Load()
	discordToken := getEnv("DISCORD_TOKEN")
	phishinToken := getEnv("PHISHIN_TOKEN")

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	phish := phishin.New(phishinToken)

	messageCreate := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		command, args, err := parseCommand(m.Content)
		if err != nil || command == UnknownCommand {
			return
		}

		switch command {
		case SetlistCommand:
			date, ok := args.(time.Time)
			if !ok {
				return
			}

			show, err := phish.ShowOnDate(context.Background(), date)
			if err != nil || show.Data.ID == 0 {
				return
			}

			setlist := setlistForShow(show.Data)
			embed, err := embedForSetlist(setlist)
			if err != nil {
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		case LastPlayedCommand:
			title, ok := args.(string)
			if !ok {
				return
			}

			lastPlayed, err := phish.LastPlayed(context.Background(), title, 4)
			if err != nil {
				return
			}

			embed, err := embedForLastPlayed(lastPlayed)
			if err != nil {
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}
	}

	dg.AddHandler(messageCreate)

	// Receive message events
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Command int

const (
	SetlistCommand Command = iota
	LastPlayedCommand
	UnknownCommand
)

// Colors: https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812
const green int = 5763719

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || len(val) == 0 {
		fmt.Println("could not find key in env:", key)
		os.Exit(1)
	}
	return val
}

func parseCommand(s string) (Command, interface{}, error) {
	r := regexp.MustCompile(`\.(setlist|lastplayed)\s+(.*)$`)
	match := r.FindStringSubmatch(s)

	if len(match) != 3 {
		return UnknownCommand, nil, nil
	}

	switch match[1] {
	case "setlist":
		t, err := dateparse.ParseAny(match[2])
		if err != nil {
			return UnknownCommand, nil, err
		}
		return SetlistCommand, t, nil
	case "lastplayed":
		return LastPlayedCommand, match[2], nil
	}

	return UnknownCommand, nil, nil
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Month(), t.Day(), t.Year())
}

func formatDuration(t time.Duration) string {
	s := t.Seconds()
	hh := int(s / 3600)
	mm := int(math.Mod(s, 3600) / 60)
	ss := int(math.Mod(s, 60))

	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

func duration(millis int) time.Duration {
	return time.Duration(millis * int(time.Millisecond))
}

type Setlist struct {
	Date     phishin.Date
	Venue    string
	Location string
	Duration time.Duration
	Sets     []*Set
}

type Set struct {
	Name     string
	Duration time.Duration
	Songs    []string
}

func setlistForShow(show phishin.Show) Setlist {
	sets := []*Set{}
	var set *Set

	for _, t := range show.Tracks {
		if set == nil {
			set = &Set{Name: t.SetName, Duration: 0 * time.Second}
		} else if set.Name != t.SetName {
			sets = append(sets, set)
			set = &Set{Name: t.SetName, Duration: 0 * time.Second}
		}
		set.Duration += duration(t.Duration)
		set.Songs = append(set.Songs, t.Title)
	}
	sets = append(sets, set)

	return Setlist{
		Date:     show.Date,
		Venue:    show.Venue.Name,
		Location: show.Venue.Location,
		Duration: duration(show.Duration),
		Sets:     sets,
	}
}

func embedForSetlist(setlist Setlist) (discordgo.MessageEmbed, error) {
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

func embedForLastPlayed(lastPlayed phishin.LastPlayed) (discordgo.MessageEmbed, error) {
	var d bytes.Buffer
	last, rest := lastPlayed.Shows[len(lastPlayed.Shows)-1], lastPlayed.Shows[:len(lastPlayed.Shows)-1]

	d.WriteString(fmt.Sprintf("It was played at %s in %s\n\n", last.Venue.Name, last.Venue.Location))
	d.WriteString("Next most recent plays 🌸:\n")
	for i := len(rest) - 1; i >= 0; i-- {
		show := rest[i]
		d.WriteString(fmt.Sprintf("🌵 %s at %s in %s\n", show.Date.Format("Monday, January 2, 2006"), show.Venue.Name, show.Venue.Location))
	}

	return discordgo.MessageEmbed{
		Color:       green,
		Title:       fmt.Sprintf("%s was last played on %s", lastPlayed.Title, last.Date.Format("Monday, January 2, 2006")),
		Description: d.String(),
	}, nil
}
