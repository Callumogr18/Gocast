# Frontend Options for Gocast

Gocast is currently a pure terminal application: `main.go` reads a menu
choice via `fmt.Scan`, calls `weathertool.FetchXML`, and the result is
printed to stdout with `%+v` (see `weather_tool/fetch_xml.go`). There is no
HTTP server, no JSON API, and no separation between "fetch data" and
"render data" yet — those two things happen in the same function.

This doc lays out realistic options for giving Gocast a nicer frontend,
from smallest effort to largest, and what each one requires from the
existing code.

## 0. Prerequisite refactor (needed for every option below)

Every option benefits from splitting `FetchXML` into two layers:

- **Data layer**: fetch + unmarshal XML into the `NationalData` /
  `ProvinceData` / `CountyForecast` structs, return `(data, error)` — no
  printing, no `fmt.Scan`.
- **Presentation layer**: takes the struct and renders it (to terminal,
  HTML, JSON, whatever).

This is a small change (a few hours) and unblocks all the options below,
including staying terminal-only with a nicer TUI.

## Option A — Nicer terminal UI (smallest change)

Keep it a CLI/TUI, no browser involved.

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** (+
  [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling,
  [Bubbles](https://github.com/charmbracelet/bubbles) for ready-made
  components like lists/spinners). Popular, well-documented, actively
  maintained by Charm. Would replace `fmt.Scan` menu with an interactive
  list/spinner while the HTTP call is in flight, and render forecast data
  in styled tables/panels.
- **[tview](https://github.com/rivo/tview)**: simpler, more "widget"
  oriented (panels, tables, forms) if you want something that looks more
  like a dashboard than a chat-style TUI.

Effort: low. Dependencies: 1-2 modules. No new architecture — still a
single Go binary.

## Option B — Local web UI, Go serves both API and HTML

Turn Gocast into a small local web server. User runs the binary, opens
`http://localhost:PORT` in a browser.

1. Add an HTTP layer (standard library `net/http`, or a light router like
   [chi](https://github.com/go-chi/chi)) with endpoints such as
   `/api/national`, `/api/province/{name}`, `/api/county/{name}` that call
   the refactored data layer and return JSON.
2. Serve a frontend from the same binary using Go's `embed.FS`
   (`//go:embed`) so the whole thing is still a single distributable
   binary with no separate deploy step:
   - **Plain HTML + a little JS** (fetch the JSON endpoints, render into
     the DOM). Zero build step, fastest to ship, fits the "simple" ask
     well.
   - **[htmx](https://htmx.org/)**: server renders HTML fragments
     directly (using Go's `html/template`), htmx swaps them into the
     page on interaction. Avoids writing a JSON API and a JS client
     separately — the Go server just returns HTML. Good fit for a small
     Go project that doesn't want a separate frontend build pipeline.
   - **A small SPA (React/Vue/Svelte + Vite)** built once with `npm run
     build` and the output embedded via `embed.FS`. More setup (Node
     toolchain, build step) but gives the most flexible/modern frontend if
     the UI is expected to grow (charts, county map, forecast history,
     etc).

Effort: medium. This is the standard way to add a "simple frontend" to a
Go CLI tool without splitting into multiple deployables.

## Option C — Separate frontend + Go API backend

Same API layer as Option B, but the frontend is a fully separate project
(e.g. a Vite/React app, or even a static site) that calls the Go server's
JSON endpoints over CORS-enabled HTTP, and is deployed independently
(Vercel/Netlify/GitHub Pages) while the Go API runs elsewhere.

Worth it only if:
- The frontend is expected to be maintained/deployed independently of the
  CLI tool, or
- You want the weather data available to something other than a browser
  UI (mobile app, other services), or
- You want frontend hot-reload / a modern JS toolchain during development
  without touching Go build tooling at all.

Effort: highest, mostly due to needing two deploy targets and CORS/auth
considerations. Overkill for a "simple frontend" unless there's a reason
to decouple.

## Recommendation

For "simple frontend" specifically: **Option B with plain HTML + fetch,
or htmx**, embedded via `embed.FS`. It keeps Gocast a single binary,
requires no Node/JS build tooling, and directly reuses the existing
struct-based data layer once it's separated from printing. Reach for
Option A instead if a browser isn't wanted at all, and Option C only if
the frontend genuinely needs to live and deploy separately from the Go
backend.
