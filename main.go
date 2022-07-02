package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
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

			embed, err := embedForLastPlayedTracks(lastPlayed)
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
	r := regexp.MustCompile(`\.(setlist|random|lastplayed|longest)(.*)$`)
	match := r.FindStringSubmatch(s)

	if len(match) == 0 {
		return UnknownCommand, nil, nil
	}

	switch match[1] {
	case "setlist":
		t, err := dateparse.ParseAny(strings.TrimSpace(match[2]))
		if err != nil {
			return UnknownCommand, nil, err
		}
		return SetlistCommand, t, nil
	case "random":
		return RandomCommand, nil, nil
	case "lastplayed":
		return LastPlayedCommand, strings.TrimSpace(match[2]), nil
	case "longest":
		return LongestCommand, strings.TrimSpace(match[2]), nil
	}

	return UnknownCommand, nil, nil
}
