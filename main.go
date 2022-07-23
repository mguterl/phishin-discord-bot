package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	_ "time/tzdata" // embed tzdata in binary

	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mguterl/phishin-discord-bot/phishin"
	"github.com/mguterl/phishin-discord-bot/phishnet"
)

func main() {
	godotenv.Load()
	discordToken := getEnv("DISCORD_TOKEN")
	phishinToken := getEnv("PHISHIN_TOKEN")
	phishNetToken := getEnv("PHISHNET_TOKEN")
	timezone := os.Getenv("TZ")

	if len(timezone) <= 0 {
		timezone = "America/Los_Angeles"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err.Error())
	}
	time.Local = loc

	p := phishnet.New(phishNetToken)
	shows, err := p.Shows(context.Background())
	if err != nil {
		fmt.Println("error fetching shows,", err)
		return
	}

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
		case RandomCommand:
			show, err := phish.RandomShow(context.Background())

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

			lastPlayed, err := phish.LastPlayedTracks(context.Background(), title, 4)
			if err != nil {
				return
			}

			embed, err := embedForLastPlayedTracks(lastPlayed, time.Now())
			if err != nil {
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		case LongestCommand:
			title, ok := args.(string)
			if !ok {
				return
			}

			longest, err := phish.LongestTracks(context.Background(), title, 3)
			if err != nil {
				return
			}

			embed, err := embedForLongestTracks(longest)
			if err != nil {
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		case DaysUntilCommand:
			date, ok := args.(time.Time)
			if !ok {
				return
			}
			message := daysUntil(time.Now(), date)

			s.ChannelMessageSend(m.ChannelID, message)
		case NextShowCommand:
			now := time.Now()
			show := phishnet.NextShow(shows.Shows, now)
			var message string

			if show == nil {
				message = "No upcoming shows"
			} else {
				difference := show.Date.Sub(now)
				days := int64(math.Ceil(difference.Hours() / 24))
				message = fmt.Sprintf("%d days until %s at the %s in %s, %s", days, formatDayOfWeek(phishin.Date(show.Date)), show.Venue, show.City, show.State)
			}

			s.ChannelMessageSend(m.ChannelID, message)
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
	RandomCommand
	LongestCommand
	DaysUntilCommand
	NextShowCommand
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

var commandRegexp = regexp.MustCompile(`\.(\w+)(.*)$`)

func parseCommand(s string) (Command, interface{}, error) {
	match := commandRegexp.FindStringSubmatch(s)

	if len(match) == 0 {
		return UnknownCommand, nil, nil
	}

	switch match[1] {
	case "random":
		return RandomCommand, nil, nil
	case "nextshow":
		return NextShowCommand, nil, nil
	case "lastplayed":
		return LastPlayedCommand, strings.TrimSpace(match[2]), nil
	case "longest":
		return LongestCommand, strings.TrimSpace(match[2]), nil
	case "setlist":
		t, err := dateparse.ParseAny(strings.TrimSpace(match[2]))
		return SetlistCommand, t, err
	case "daysuntil":
		t, err := dateparse.ParseAny(strings.TrimSpace(match[2]))
		return DaysUntilCommand, t, err
	}

	return UnknownCommand, nil, nil
}
