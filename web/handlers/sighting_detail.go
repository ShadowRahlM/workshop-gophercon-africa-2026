package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"uuid"
	"workshop/web/views/pages"
)

func SightingDetail(logger *slog.Logger, sightings SightingsFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract UUID from path: /sightings/{id}
		idStr := strings.TrimPrefix(r.URL.Path, "/sightings/")
		if idStr == "" {
			renderError(ctx, w, logger, http.StatusNotFound, "Sighting ID is required.")
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			renderError(ctx, w, logger, http.StatusNotFound, "Invalid sighting ID.")
			return
		}

		sight, err := sightings.Find(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				renderError(ctx, w, logger, http.StatusNotFound, "Sighting not found.")
				return
			}
			logger.ErrorContext(ctx, "handlers: finding sighting: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sighting.")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.SightingDetail(sight).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering sighting detail: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
