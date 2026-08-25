# Gocast: feature ideas to make it feel like a real product

## Context

Gocast is a Go weather app for Ireland with two thin UIs (Bubble Tea CLI, htmx web UI) over a narrow data layer (`weather_tool/`) that scrapes three met.ie XML feeds: national summary text, provincial forecast+pollen, and a 3–4 day county table. There's no DB, no cache, no auth, no router beyond stdlib `ServeMux`, no tests, no CI/Docker. A previous agent session already left a technical-debt/refactor backlog in `docs/agent/HANDOFF.md` (typed locations, TTL cache, `embed.FS`, middleware, tests) — that's infrastructure cleanup, not new product surface, and is treated as separate from this list.

This file is the output of research into what additions would make Gocast more usable and closer to a real product — a prioritized list of concrete, grounded feature ideas, each with rationale, rough approach, and effort/dependency notes.

## Research finding worth flagging up front

Met Éireann publishes a **weather warnings** feed that Gocast doesn't use at all today: `https://www.met.ie/warningsxml/rss.xml` (colour-coded national/yellow/orange/red warnings, county-taggable). Gocast currently has zero severe-weather awareness — for a "real" weather product, warnings are arguably more valuable than the forecast text it already has. This is the single highest-leverage net-new feature.

Also confirmed: met.ie's XML includes a `<symbol id="PartlyCloud" code="partlycloudy_day">`-style icon code per forecast entry — Gocast currently only surfaces raw weather text, so icons are a cheap visual upgrade once that field is parsed (may need a struct field added, similar to existing `Day` fields in `structs_formats.go`).

A public rain-radar image/embed URL could **not** be confirmed as stable/current (only an old, likely-dead `Web_radar.gif` link turned up) — not recommending it as a v1 item; would need direct confirmation from met.ie before relying on it.

## Recommended additions, roughly prioritized

**Tier 1 — no new dependencies, no persistence, high visible impact**

1. **Weather warnings banner** — fetch `https://www.met.ie/warningsxml/rss.xml`, parse (stdlib `encoding/xml`, same pattern as `fetch_xml.go`), show a colour-coded alert banner on the home page and/or per-county view when a warning is active. New file e.g. `weather_tool/fetch_warnings.go` mirroring existing fetch functions; new handler + template fragment.
2. **Weather icons** — parse the `symbol`/`code` attribute already present in the XML (currently dropped), map codes to a small icon set (emoji or a handful of inline SVGs — no image hosting needed), render next to each forecast day/summary.
3. **Bookmarkable/shareable URLs** — currently county/province selection is pure htmx client interaction with no URL state, so a link to "Dublin's forecast" can't be shared. Add query params (`/county?name=dublin`) that both drive the fetch and get reflected via `hx-push-url="true"`, so the back button and shared links work. Pure frontend/routing change, no new deps.
4. **"Remember my county" via cookie or localStorage** — auto-select the user's last-viewed county/province on return visits. No DB needed; a simple cookie set from the handler or `localStorage` + a tiny inline script is enough for a single-user-per-browser feel.
5. **Loading/error UX polish** — HANDOFF/UI-ENHANCEMENTS docs suggest most of this is done (spinners, error alert box); worth a quick pass to confirm consistency across all four routes rather than treating as new work.

**Tier 2 — still no backend dependencies, more implementation effort**

6. **County comparison view** — pick two counties, render forecasts side-by-side. Reuses existing `FetchCounty` twice; new template/route only.
7. **Min/max temperature trend strip** — the county forecast already returns `MinTemp`/`MaxTemp` per day; render a tiny inline SVG sparkline/bar strip across the fetched days instead of only a table. No charting library needed at this data volume — can hand-roll simple bar/line SVG.
8. **PWA basics** — a manifest.json + minimal service worker to cache the last successfully rendered page for offline viewing, and make the site "installable" on mobile. No backend change; makes the app feel materially more like a real product on a phone.

**Tier 3 — needs new dependencies and/or persistence (bigger lift, flag before starting)**

9. **Saved favourite locations across visits/devices + optional email/push digest of warnings** — needs a lightweight DB (SQLite would be enough, no server needed) and, for notifications, a mail-sending dependency or a browser Push API + service worker. This is the natural "grows into a real product" feature but is the biggest scope jump from the current zero-persistence codebase.
10. **User accounts/auth** — only worth it if favourites/notifications (item 9) are wanted across devices; otherwise skip, since cookie/localStorage covers the single-device case in Tier 1.

## Suggested next step

Tier 1 items are all cheap, additive, and don't conflict with the existing HANDOFF.md refactor backlog — they could be picked up in any order. The warnings banner (#1) is the standout recommendation: it's genuinely missing functionality (not just polish) and uses a real, confirmed met.ie feed.

## Verification approach (once any item is implemented)

- Run `go build ./...` and the existing manual flow (`go run .` → `y` for web UI) to click through the affected route(s) in a browser.
- For new XML parsing (warnings, icon codes), fetch the live feed once via `curl` to confirm the real shape matches the struct before wiring it into a handler — met.ie's docs have previously been stale (README already documents this mismatch for the existing feeds).
- No test suite exists yet; if the user wants confidence beyond manual clicking, pair any Tier 1/2 item with a small `_test.go` for the new parsing/handler logic (table-driven, no live network call) rather than expanding scope elsewhere.
