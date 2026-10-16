package explore

import (
	"cmp"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"workshop/cmd/themes/internal/layout"
	"workshop/cmd/themes/internal/themeutil"
	"workshop/internal/domain/sighting"
)

// RecordProps is everything one sighting's own page renders.
const fonts = "https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,700;0,900;1,400&family=Inter:wght@300;400;500;600&display=swap"

type RecordProps struct {
	Sighting sighting.Sighting
	Now      time.Time
	Thumb    string
	Map      *layout.MapProps
}

// StatisticsProps carries the headline figures and the three charts.
type StatisticsProps struct {
	Summary   sighting.Summary
	Species   []Bar
	Countries []Bar
	Years     []Bar
}

// GalleryProps carries one photograph for each sighting that has one.
type GalleryProps struct {
	Photos  []Photo
	Summary sighting.Summary
}

// AtlasProps places every sighting on the map, with the list beneath it.
type AtlasProps struct {
	Sightings []sighting.Sighting
	Map       *layout.MapProps
}

// Photo is the first still image filed against a sighting.
type Photo struct {
	Sighting sighting.Sighting
	URL      string
}

// Bar is one row of a chart, scaled against the largest row in its set.
type Bar struct {
	Label   string
	Count   int
	Percent int
}

// NewRecordProps arranges one sighting for its own page.
func NewRecordProps(s sighting.Sighting, now time.Time) RecordProps {
	return RecordProps{
		Sighting: s,
		Now:      now,
		Thumb:    themeutil.Thumbnail(s),
		Map:      mapOf([]sighting.Sighting{s}),
	}
}

// recordPath is where a sighting's own page lives.
func recordPath(s sighting.Sighting) string {
	return "/sightings/" + s.ID.String()
}

// NewStatisticsProps works out the figures and charts from the sightings.
func NewStatisticsProps(sightings []sighting.Sighting) StatisticsProps {
	summary := sighting.Summarise(sightings)

	return StatisticsProps{
		Summary:   summary,
		Species:   bars(sighting.Top(sighting.CountBy(sightings, speciesName), 10)),
		Countries: bars(sighting.Top(sighting.CountBy(sightings, countryName), 10)),
		Years:     yearBars(sightings, summary),
	}
}

// NewGalleryProps keeps the sightings that came with a photograph.
func NewGalleryProps(sightings []sighting.Sighting) GalleryProps {
	photos := make([]Photo, 0, len(sightings))
	for _, s := range sightings {
		if url := themeutil.Thumbnail(s); url != "" {
			photos = append(photos, Photo{Sighting: s, URL: url})
		}
	}

	return GalleryProps{Photos: photos, Summary: sighting.Summarise(sightings)}
}

// NewAtlasProps places the sightings on the map.
func NewAtlasProps(sightings []sighting.Sighting) AtlasProps {
	return AtlasProps{Sightings: sightings, Map: mapOf(sightings)}
}

// mapOf returns the map for the sightings, or nil when none has a position.
func mapOf(sightings []sighting.Sighting) *layout.MapProps {
	points := themeutil.MapPoints(sightings)
	if len(points) == 0 {
		return nil
	}
	return &layout.MapProps{Points: points}
}

// bars scales a leaderboard against its first row, which is the largest.
func bars(tallies []sighting.Tally) []Bar {
	top := 0
	if len(tallies) > 0 {
		top = tallies[0].Count
	}

	out := make([]Bar, 0, len(tallies))
	for _, t := range tallies {
		out = append(out, Bar{Label: t.Name, Count: t.Count, Percent: themeutil.Percent(t.Count, top)})
	}
	return out
}

// yearBars keeps every year in the span, so quiet years still show as gaps.
func yearBars(sightings []sighting.Sighting, summary sighting.Summary) []Bar {
	counts := sighting.CountByYear(sightings)

	top := 0
	for _, n := range counts {
		top = max(top, n)
	}

	years := themeutil.Years(summary)
	out := make([]Bar, 0, len(years))
	for _, y := range years {
		out = append(out, Bar{
			Label:   strconv.Itoa(y),
			Count:   counts[y],
			Percent: themeutil.Percent(counts[y], top),
		})
	}
	return out
}

// speciesName labels a tally row with the name the pages show.
func speciesName(s sighting.Sighting) string {
	return themeutil.DisplayName(s.Species)
}

// countryName labels a tally row with the country a sighting was made in.
func countryName(s sighting.Sighting) string {
	return s.Location.Country
}

// basis turns GBIF's shouted enum into something a reader can take in.
func basis(s sighting.Sighting) string {
	if s.BasisOfRecord == "" {
		return "Not stated"
	}
	words := strings.ReplaceAll(strings.ToLower(s.BasisOfRecord), "_", " ")
	return strings.ToUpper(words[:1]) + words[1:]
}

// count renders a whole number for the templates.
func count(n int) string {
	return strconv.Itoa(n)
}

// condition names the present weather, or says it was not recorded.
func condition(description string) string {
	return cmp.Or(description, "Conditions not recorded")
}

// when renders the moment a weather reading was taken.
func when(t time.Time) string {
	return t.Format("2 January 2006, 15:04")
}

// celsius renders a temperature in degrees Celsius.
func celsius(v float64) string {
	return fmt.Sprintf("%.1f °C", v)
}

// millimetres renders rainfall or other precipitation.
func millimetres(v float64) string {
	return fmt.Sprintf("%.1f mm", v)
}

// kmh renders a wind speed.
func kmh(v float64) string {
	return fmt.Sprintf("%.0f km/h", v)
}

// percentage renders a cover or humidity figure.
func percentage(n int) string {
	return fmt.Sprintf("%d%%", n)
}

// barWidth sets the fill of a chart bar, as a percentage of its track.
func barWidth(percent int) string {
	return fmt.Sprintf("width: %d%%", percent)
}

var compassPoints = [...]string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}

// compass names the direction a wind blew from, given in degrees.
func compass(degrees int) string {
	d := ((degrees % 360) + 360) % 360
	return compassPoints[int(math.Round(float64(d)/22.5))%len(compassPoints)]
}
