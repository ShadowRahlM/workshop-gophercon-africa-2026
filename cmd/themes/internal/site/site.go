package site

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"time"
	"uuid"
	"workshop/cmd/themes/internal/antarctica"
	"workshop/cmd/themes/internal/humpback"
	"workshop/cmd/themes/internal/kenya"
	"workshop/cmd/themes/internal/uganda"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/internal/domain/species"

	"github.com/a-h/templ"
)

// assets holds the files served under /assets.
// Note: The embedded paths keep their "assets/" prefix.
//
//go:embed assets
var assets embed.FS

// Lister is all the themes need from the sighting domain. Declaring it
// here keeps them reading the starter's data without depending on how
// the starter serves it.
type Lister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
	Find(ctx context.Context, id uuid.UUID) (sighting.Sighting, error)
}

// Handler serves the themed pages.
type Handler struct {
	logger    *slog.Logger
	sightings Lister
}

// New constructs a Handler over anything that can list sightings.
func New(logger *slog.Logger, sightings Lister) *Handler {
	return &Handler{
		logger:    logger,
		sightings: sightings,
	}
}

// page renders one theme's markup from the sightings its filter matched.
type page func([]sighting.Sighting, time.Time) templ.Component

// Routes registers a route per themed page and the assets they load.
// The links between the pages live in the layout package.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /assets/", http.FileServerFS(assets))

	// The workshop is in Nairobi, so the front door opens on Kenya.
	mux.Handle("GET /{$}", http.RedirectHandler("/kenya", http.StatusFound))

	mux.HandleFunc("GET /antarctica", h.theme(
		sighting.ListFilter{Country: country.Antarctica},
		func(s []sighting.Sighting, now time.Time) templ.Component {
			return antarctica.Page(antarctica.NewProps(s, now))
		},
	))

	mux.HandleFunc("GET /kenya", h.theme(
		sighting.ListFilter{Country: country.Kenya},
		func(s []sighting.Sighting, now time.Time) templ.Component {
			return kenya.Page(kenya.NewProps(s, now))
		},
	))

	mux.HandleFunc("GET /humpback", h.theme(
		sighting.ListFilter{Animals: []species.Animal{species.AnimalHumpbackWhale}},
		func(s []sighting.Sighting, now time.Time) templ.Component {
			return humpback.Page(humpback.NewProps(s, now))
		},
	))

	mux.HandleFunc("GET /uganda", h.theme(
		sighting.ListFilter{Country: country.Uganda},
		func(s []sighting.Sighting, now time.Time) templ.Component {
			return uganda.Page(uganda.NewProps(s, now))
		},
	))
}

// theme builds the handler for one themed page: list with the filter,
// hand the result to the page, render. The themes differ only in those
// two values, so the plumbing is written once.
func (h *Handler) theme(filter sighting.ListFilter, render page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sightings, err := h.sightings.List(ctx, filter)
		if err != nil {
			h.logger.ErrorContext(ctx, "site: listing sightings: "+err.Error())
			http.Error(w, "Unable to load sightings.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := render(sightings, time.Now()).Render(ctx, w); err != nil {
			// The status is already written, so this can only be logged.
			h.logger.ErrorContext(ctx, "site: rendering theme: "+err.Error())
		}
	}
}
