# Gocast — Handoff / Context Doc

Purpose-built for picking this project back up in a future session. Keep it
concise and update it when the shape of the project changes materially
(new layer, new dependency, architecture shift) — not on every commit.

## What it is

Gocast is a Go CLI/TUI that fetches Irish weather forecasts (national,
regional/province, county) and prints them to the terminal. Data source is
XML feeds from `met.ie` (`www.met.ie/Open_Data/xml/...`), fetched directly
over HTTP and unmarshaled into Go structs.

**Note:** `README.md` describes the source as `https://weather.apis.ie/docs`
and `docs/README-XML.md` documents a GraphQL request/response shape, but
the actual implementation in `weather_tool/fetch_xml.go` hits met.ie XML
endpoints directly with `encoding/xml`. The docs predate/describe a
different API than what's implemented — treat `docs/README-XML.md` as
historical reference for the data shape, not as ground truth for the code.

## Layout

- `main.go` — entrypoint; calls `display.RunSelector()` for the menu choice,
  then `weathertool.FetchXML(choice)`.
- `display/cli.go` — Bubble Tea (v2) TUI: a simple list selector
  (National / Regional / County) using `charm.land/bubbletea`,
  `charm.land/bubbles`, `charm.land/lipgloss`.
- `weather_tool/fetch_xml.go` — HTTP fetch + XML unmarshal + print, per
  forecast type. Location input (province/county) is read via `fmt.Scan`
  *after* the TUI selector exits (not part of the Bubble Tea flow).
- `weather_tool/structs_formats.go` — XML-tagged structs: `NationalData`,
  `ProvinceData`, `CountyForecast`/`CountyData`/`DayResult`.
- `docs/README-XML.md` — GraphQL-shaped request/response reference (see
  discrepancy note above).
- `docs/FRONTEND-OPTIONS.md` — options doc for eventually giving Gocast a
  web frontend (Options A/B/C, recommends embedding a plain HTML/htmx UI
  via `embed.FS` if that's ever pursued). Prerequisite for any of them:
  split `FetchXML` into a data layer (fetch+unmarshal, no printing) and a
  presentation layer — currently they're fused.

## Current state (as of 2026-08-22)

- Branch `develop` is even with `origin/develop`; `master` is the PR-merge
  target branch.
- Recent work replaced the old `fmt.Scan`-only menu with a Bubble Tea list
  selector (`display/cli.go`) — see commits `58b1391`, `a9b34cc`, `7ee58cf`.
- `docs/FRONTEND_OPTIONS.md` (underscore) was renamed to
  `docs/FRONTEND-OPTIONS.md` (hyphen); the old path showed as deleted in
  `git status` until that rename was committed.
- No tests exist yet (`go test ./...` finds nothing to run).
- No HTTP server / API layer — still a single-binary CLI.

## Conventions / practices to follow

- **Module path:** `github.com/Callumogr18/Gocast`, Go 1.26.1.
- **Charm libraries are on `v2` and the `charm.land/...` import path**, not
  the older `github.com/charmbracelet/...` path — don't mix them when
  adding new Bubble Tea/Bubbles/Lipgloss code.
- **Fetch/parse/print are currently one function per forecast type** in
  `fetch_xml.go`. If you extend this, prefer splitting fetch+unmarshal from
  printing rather than adding a fourth mixed-concern branch — this is the
  prerequisite refactor `docs/FRONTEND-OPTIONS.md` already calls out.
- **Errors are printed and swallowed** (`fmt.Println("failed...", err); return`)
  rather than returned — consistent with a CLI-only tool today, but will
  need to change to real `error` returns if a data layer is split out for
  an API/web frontend.
- Branch flow so far: feature branches → PR into `develop` → PR into
  `master`. No CI configured yet.

## Suggested next steps (not started)

1. Do the fetch/print split described in `docs/FRONTEND-OPTIONS.md` §0 —
   unblocks tests and any frontend option.
2. Add basic unit tests around XML unmarshaling (the structs are the part
   most likely to silently break if met.ie changes its XML shape).
3. Decide whether `docs/README-XML.md` should be rewritten to match the
   real met.ie XML feeds, or removed if it's no longer useful as reference.
