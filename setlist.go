package main

import (
	"time"

	"github.com/mguterl/phishin-discord-bot/phishin"
)

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

func setlistForShow(show phishin.Show) *Setlist {
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

	return &Setlist{
		Date:     show.Date,
		Venue:    show.Venue.Name,
		Location: show.Venue.Location,
		Duration: duration(show.Duration),
		Sets:     sets,
	}
}

func duration(millis int) time.Duration {
	return time.Duration(millis * int(time.Millisecond))
}
