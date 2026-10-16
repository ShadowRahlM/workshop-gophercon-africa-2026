# Wildlife Sightings — Design

## Theme: Night Safari

A dark, immersive design inspired by night wildlife photography — deep backgrounds with warm gold accents that highlight the animals and their stories.

## Palette (CSS Custom Properties)

| Token | Hex | Use |
|-------|-----|-----|
| `--color-bg` | `#0C0A09` | Page background |
| `--color-surface` | `#1C1917` | Cards, sections |
| `--color-surface-raised` | `#292524` | Elevated elements |
| `--color-border` | `#44403C` | Borders, dividers |
| `--color-text` | `#FAFAFA` | Primary text |
| `--color-text-muted` | `#A8A29E` | Secondary text |
| `--color-text-dim` | `#78716C` | Tertiary text, labels |
| `--color-gold` | `#D4A017` | Accents, links, stats |
| `--color-gold-light` | `#EBC017` | Hover states |
| `--color-gold-dim` | `#B8860B` | Icons, subtle accents |

## Typography

- **Headings**: Outfit (400, 500, 600, 700) — geometric and modern
- **Body**: Work Sans (400, 500, 600) — clean and readable

Loaded from Google Fonts via preconnect links in the layout.

## Spacing & Radii

- `--radius`: 12px (cards, hero section)
- `--radius-sm`: 8px (buttons, small elements)
- Container max-width: 1200px
- Card gap: 24px

## Components

### Header
- Fixed at top with logo (binoculars SVG) and site name
- Navigation links with hover background
- Border-bottom separator

### Hero Section
- Gradient background with gold radial glow
- Large title with binoculars icon
- Subtitle with date range
- Four stat columns: Sightings, Species, Countries, Photos
- Gold stat values with uppercase labels

### Sparkline Chart
- SVG line chart showing sightings over time
- Gold line with subtle fill underneath
- Accessible aria-label with peak year info

### Sighting Cards
- Dark surface with border, hover lifts with shadow
- Two-column grid: gallery + info
- Gold icon accents in metadata
- Weather summary in raised surface box
- Gallery images with hover zoom effect

### Map
- Dark tiles via `map--dark` class (inverted Leaflet tiles)
- Gold markers with dark stroke
- Rounded corners with border

### Statistics Page
- Stat cards with gold values and hover glow
- Tables with gold count column
- Year chart bars with gold gradient

### Error Page
- Gold button with dark text
- Dim icon and heading

## Accessibility

- All text meets 4.5:1 contrast ratio on dark backgrounds
- Focus-visible outlines in gold
- `prefers-reduced-motion` disables all transitions and transforms
- Screen reader labels on charts and map
- Semantic HTML structure maintained

## Responsive

- Single column layout below 640px
- Hero stats wrap to multiple rows
- Cards stack vertically on mobile
- Container padding reduces on small screens
