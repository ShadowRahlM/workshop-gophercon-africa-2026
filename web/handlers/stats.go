package handlers

import (
	"log/slog"
	"net/http"
	"sort"
	"workshop/internal/domain/sighting"
	"workshop/web/views/pages"
)

func Stats(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, sighting.ListFilter{
			Limit: 500,
		})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings for stats: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		summary := sighting.Summarise(items)
		topSpecies := sighting.Top(sighting.CountBy(items, func(s sighting.Sighting) string {
			return s.Species.ScientificName
		}), 10)
		topCountries := sighting.Top(sighting.CountBy(items, func(s sighting.Sighting) string {
			return s.Location.Country
		}), 10)
		byYear := sighting.CountByYear(items)

		// Convert year map to sorted rows with percentages
		yearRows := buildYearRows(byYear)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Stats(summary, topSpecies, topCountries, yearRows).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering stats: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}

func buildYearRows(byYear map[int]int) []pages.YearRow {
	// Sort years
	years := make([]int, 0, len(byYear))
	for year := range byYear {
		years = append(years, year)
	}
	sort.Ints(years)

	// Find max for scaling
	maxCount := 0
	for _, count := range byYear {
		if count > maxCount {
			maxCount = count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	// Build rows
	rows := make([]pages.YearRow, 0, len(years))
	for _, year := range years {
		count := byYear[year]
		percent := int(float64(count) / float64(maxCount) * 100)
		rows = append(rows, pages.YearRow{
			Year:    year,
			Count:   count,
			Percent: percent,
		})
	}
	return rows
}
