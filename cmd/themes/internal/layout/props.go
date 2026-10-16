package layout

import "github.com/a-h/templ"

// BaseProps is everything the document shell needs to know about the
// page inside it. Only Title is required.
type BaseProps struct {
	Title       string
	Description string

	// Class goes on <body>, MainClass on <main>.
	Class     string
	MainClass string

	// Fonts is a webfont stylesheet URL, preconnect hints included.
	Fonts string

	// Stylesheet is the page's own CSS, loaded last so it wins.
	Stylesheet string

	// Skip is the shortcut past the header, for anyone who cannot skim
	// their way down the page. The zero value renders no link.
	Skip SkipLink

	// Header and Footer sit outside <main>, so they cannot come through
	// the children slot.
	Header templ.Component
	Footer templ.Component

	// Map, when set, loads Leaflet and writes its markers out as JSON.
	// How the map looks stays with the theme, on the map element.
	Map *MapProps

	// Scripts, when set, is rendered just before </body>.
	Scripts templ.Component
}

// SkipLink points at an id, written without the leading #.
type SkipLink struct {
	Target string
	Label  string
}

// MapProps carries the markers themeutil.MapPoints built.
type MapProps struct {
	Points any
}
