package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"uuid"
	"workshop/internal/domain/sighting"
	"workshop/web/views/pages"
)

type SightingsLister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
}

type SightingsFinder interface {
	Find(ctx context.Context, id uuid.UUID) (sighting.Sighting, error)
}

func Home(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, sighting.ListFilter{
			Limit: 50,
		})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		summary := sighting.Summarise(items)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Home(items, summary).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
