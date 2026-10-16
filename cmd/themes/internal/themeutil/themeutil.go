package themeutil

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
)

// MapPoint is one marker handed to Leaflet as JSON.
type MapPoint struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Title  string  `json:"title"`
	Date   string  `json:"date"`
	Place  string  `json:"place"`
	Count  string  `json:"count,omitzero"`
	Image  string  `json:"image,omitzero"`
	Anchor string  `json:"anchor"`
}

// DisplayName prefers the name a reader would recognise, falling back
// through the taxonomy as GBIF records thin out.
func DisplayName(sp species.Species) string {
	return cmp.Or(sp.VernacularName, sp.CanonicalName, sp.ScientificName, "Unknown species")
}

// Place joins whichever parts of the location GBIF supplied.
func Place(l sighting.Location) string {
	parts := slices.DeleteFunc(
		[]string{l.Locality, l.StateProvince, l.Country},
		func(s string) bool { return s == "" },
	)
	return cmp.Or(strings.Join(parts, ", "), "Location unknown")
}

// Locale is the shorter form of Place, for cards that already show the
// country somewhere else.
func Locale(l sighting.Location) string {
	return cmp.Or(l.Locality, l.StateProvince, l.Country, "Unknown")
}

// Coordinates renders a position the way a field record would, with a
// hemisphere letter rather than a minus sign.
func Coordinates(l sighting.Location) string {
	return fmt.Sprintf("%.3f%s, %.3f%s",
		math.Abs(l.Latitude), hemisphere(l.Latitude, "N", "S"),
		math.Abs(l.Longitude), hemisphere(l.Longitude, "E", "W"),
	)
}

// Individuals renders how many animals were seen in a single sighting.
func Individuals(n *int) string {
	if n == nil {
		return ""
	}
	if *n == 1 {
		return "1 individual"
	}
	return fmt.Sprintf("%d individuals", *n)
}

// Recorder names whoever filed the record, for the many GBIF rows that
// leave the field blank.
func Recorder(s sighting.Sighting) string {
	return cmp.Or(s.RecordedBy, "Unattributed")
}

// Thumbnail returns the first still image on a sighting, or an empty
// string when it has none.
func Thumbnail(s sighting.Sighting) string {
	for _, m := range s.Media {
		if m.Type == "StillImage" && m.URL != "" {
			return m.URL
		}
	}
	return ""
}

// Credit names whoever took the first photograph on a sighting, with
// the licence GBIF gave it when there is one.
func Credit(s sighting.Sighting) string {
	for _, m := range s.Media {
		if m.Type != "StillImage" || m.URL == "" {
			continue
		}
		name := cmp.Or(m.Creator, m.RightsHolder, "Unattributed")
		if m.License != "" {
			return name + " · " + Licence(m.License)
		}
		return name
	}
	return ""
}

// Licence shortens the licence URLs GBIF returns to their name.
func Licence(raw string) string {
	switch {
	case strings.Contains(raw, "by-nc-sa"):
		return "CC BY-NC-SA"
	case strings.Contains(raw, "by-nc-nd"):
		return "CC BY-NC-ND"
	case strings.Contains(raw, "by-nc"):
		return "CC BY-NC"
	case strings.Contains(raw, "by-sa"):
		return "CC BY-SA"
	case strings.Contains(raw, "by-nd"):
		return "CC BY-ND"
	case strings.Contains(raw, "publicdomain"), strings.Contains(raw, "zero"):
		return "CC0"
	case strings.Contains(raw, "creativecommons.org/licenses/by"):
		return "CC BY"
	}
	return raw
}

// AltText describes a sighting's image for anyone who cannot see it.
func AltText(s sighting.Sighting) string {
	return fmt.Sprintf("%s photographed at %s on %s",
		DisplayName(s.Species), Place(s.Location), s.HappenedAt.Format("2 January 2006"))
}

// Anchor is a stable element id for a sighting, so a map marker can
// link to the record's row further down the page.
func Anchor(s sighting.Sighting) string {
	return "sighting-" + s.ID.String()
}

// Percent scales count against total, for bars drawn in CSS. It returns
// a whole number between 0 and 100.
func Percent(count, total int) int {
	if total <= 0 {
		return 0
	}
	return int(math.Round(float64(count) / float64(total) * 100))
}

// Span describes the period a set of sightings covers.
func Span(s sighting.Summary) string {
	if s.Sightings == 0 {
		return "No records"
	}
	if s.Earliest.Year() == s.Latest.Year() {
		return fmt.Sprint(s.Earliest.Year())
	}
	return fmt.Sprintf("%d–%d", s.Earliest.Year(), s.Latest.Year())
}

// Years lists every year between the earliest and latest sighting, so a
// chart keeps its gaps instead of closing them up.
func Years(s sighting.Summary) []int {
	if s.Sightings == 0 {
		return nil
	}
	years := make([]int, 0, s.Latest.Year()-s.Earliest.Year()+1)
	for y := s.Earliest.Year(); y <= s.Latest.Year(); y++ {
		years = append(years, y)
	}
	return years
}

// MapPoints turns sightings into the markers the map script reads.
func MapPoints(sightings []sighting.Sighting) []MapPoint {
	points := make([]MapPoint, 0, len(sightings))
	for _, s := range sightings {
		if s.Location.Latitude == 0 && s.Location.Longitude == 0 {
			continue
		}
		points = append(points, MapPoint{
			Lat:    s.Location.Latitude,
			Lng:    s.Location.Longitude,
			Title:  DisplayName(s.Species),
			Date:   s.HappenedAt.Format("2 Jan 2006"),
			Place:  Place(s.Location),
			Count:  Individuals(s.IndividualCount),
			Image:  Thumbnail(s),
			Anchor: Anchor(s),
		})
	}
	return points
}

// Since renders how long ago a sighting happened, relative to now.
func Since(t time.Time, now time.Time) string {
	days := int(now.Sub(t).Hours() / 24)
	switch {
	case days < 0:
		return "Upcoming"
	case days == 0:
		return "Today"
	case days == 1:
		return "Yesterday"
	case days < 30:
		return fmt.Sprintf("%d days ago", days)
	case days < 365:
		return fmt.Sprintf("%d months ago", days/30)
	case days < 730:
		return "A year ago"
	}
	return fmt.Sprintf("%d years ago", days/365)
}

// hemisphere picks the letter for a signed coordinate.
func hemisphere(v float64, positive, negative string) string {
	if v < 0 {
		return negative
	}
	return positive
}
