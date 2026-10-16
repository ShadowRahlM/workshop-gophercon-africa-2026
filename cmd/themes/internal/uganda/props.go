package uganda

import (
	"strconv"
	"strings"
	"time"
	"workshop/cmd/themes/internal/themeutil"
	"workshop/internal/domain/sighting"
)

// Props is everything the Uganda wildlife page renders.
type Props struct {
	Entries []Entry
	Summary sighting.Summary
	Species []sighting.Tally
	Regions []sighting.Tally
	Points  []themeutil.MapPoint
	Now     time.Time
}

// Entry is one sighting written up as a numbered article.
type Entry struct {
	Sighting sighting.Sighting
	Number   int
}

// NewProps arranges the sightings into a page, newest first.
func NewProps(sightings []sighting.Sighting, now time.Time) Props {
	entries := make([]Entry, len(sightings))
	for i, s := range sightings {
		entries[i] = Entry{Sighting: s, Number: i + 1}
	}

	regions := make(map[string]int)
	for _, s := range sightings {
		if name := s.Location.StateProvince; name != "" {
			regions[name]++
		}
	}

	return Props{
		Entries: entries,
		Summary: sighting.Summarise(sightings),
		Species: sighting.Top(sighting.CountBy(sightings, speciesName), 5),
		Regions: sighting.Top(regions, 5),
		Points:  themeutil.MapPoints(sightings),
		Now:     now,
	}
}

// folio renders an entry's number with zero-padding.
func folio(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// count renders a tally for the templates.
func count(n int) string {
	return strconv.Itoa(n)
}

// standfirst is the opening sentence, which changes with how much there is.
func standfirst(s sighting.Summary) string {
	switch s.Sightings {
	case 0:
		return "No sightings have been recorded in Uganda yet."
	case 1:
		return "A single record from the Pearl of Africa, in full."
	}
	return "Every wildlife sighting logged across Uganda, read end to end."
}

// basis turns GBIF's shouted enum into something a reader can take in.
func basis(s sighting.Sighting) string {
	if s.BasisOfRecord == "" {
		return "Not stated"
	}
	words := strings.ReplaceAll(strings.ToLower(s.BasisOfRecord), "_", " ")
	return strings.ToUpper(words[:1]) + words[1:]
}

// speciesName labels a tally row with the name this theme shows.
func speciesName(s sighting.Sighting) string {
	return themeutil.DisplayName(s.Species)
}
