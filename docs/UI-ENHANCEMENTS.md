# UI Enhancements

Known issues in the HTMX-based web UI (`htmx/html/*.html`, `htmx/server/*.go`) and their fixes.

## 1. Escape characters shown literally in forecast text

**Resolved.** Fixed via `cleanString` in `weather_tool/fetch_helpers.go`,
called from `fetchAndRead` — literal `\n` sequences are stripped/collapsed
before the data ever reaches a template.


**What happens:** Free-text forecast fields (`Outlook`, `Today`, `Tomorrow`, `Pollen`) display raw `\n\n` characters inline instead of line breaks. For example, the National page renders text like:

> Overview: Unsettled with rain...\n\nWednesday night: Showers...\n\nThursday will st...

**Why:** These fields are unmarshalled straight from met.ie's XML feed as plain `string` values (`weather_tool/structs_formats.go`, via `xml.Unmarshal` in `weather_tool/fetch_xml.go`) with no cleanup step. The upstream feed contains literal backslash-n sequences (not real newline bytes). They're then inserted into the templates (e.g. `htmx/html/national.html:4`) via `html/template`, which HTML-escapes markup characters like `<`, `>`, `&`, `"`, `'` but has no concept of `\n` as a line break — it passes the two literal characters straight through to the browser. `Day.WeatherLabel()` (`weather_tool/structs_formats.go:52-55`) already does a similar cleanup (replacing `_` with a space) for a different field, but no equivalent exists for these free-text fields.

**Fix:** Add a small sanitizing helper, e.g. `cleanText(s string) string`, that replaces literal `\n` sequences with real line breaks. Apply it either when populating the struct fields after parsing, or via a template `FuncMap` at render time. Since converting to actual `<br>` tags requires marking the output as safe HTML (`template.HTML`), the simplest safe approach is to split on `\n` and render each segment as its own `<p>` element rather than injecting raw HTML.

## 2. Confirm dialog on National is inconsistent and unnecessary

**What happens:** The National button (`htmx/html/home.html:16-27`) is wrapped in a SweetAlert2 confirmation dialog — clicking it fires `Swal.fire(...)`, and only on confirm does it call `htmx.trigger(this, 'confirmed')`, which the button listens for via `hx-trigger="confirmed"` (not a normal click). A plain click does nothing until the modal is dismissed. Regional (`home.html:30-36`) and County (`home.html:40-74`) use plain `<select>` elements with `hx-trigger="change"` and fire immediately, with no confirmation step.

**Why:** `nationalHandler` (`htmx/server/handlers.go:46-54`) takes no query parameters and performs no destructive or ambiguous action — it's a simple GET, functionally simpler than Regional/County which take user-supplied `province`/`county` params yet have no confirmation gate at all. There's no evident reason National alone needs an "are you sure?" step; it reads as leftover/inconsistent UX rather than an intentional safeguard. The button's label ("Click Me") also doesn't describe the action.

**Fix:** Remove the SweetAlert2 confirm dialog and the `hx-trigger="confirmed"` binding for National. Make it a plain button with `hx-get="/national"` using the default click trigger, matching Regional/County behavior. Rename the button label to something descriptive, e.g. "National Forecast".

## 3. No styling — unstyled default browser layout

**Resolved.** See `docs/agent/HANDOFF.md` "UI design pass" —
`htmx/static/style.css` + a `GET /static/` route + classed templates.


**What happens:** There is no CSS anywhere in the project — no `.css` files, no `<style>` blocks, no `style=` attributes in any template. The only external assets are two CDN `<script>` tags (htmx, SweetAlert2) loaded in `home.html`. All pages render with default browser styling for headings, tables, and selects, with no responsive layout and no visual separation between the National/Regional/County sections beyond plain `<h2>` headers.

**Why:** The UI was built to validate HTMX wiring and data flow first; no styling pass has been done yet. There's also no static file route registered in `htmx/server/http_server.go`, so there's currently no mechanism to serve a CSS file even if one were added — and the fragment templates (`national.html`, `regional.html`, `county.html`) have no `<head>` of their own since they're rendered as HTMX-swapped fragments into `home.html`.

**Fix:**
- Add a `static/` directory containing a `style.css`, and register a static file handler (`http.FileServer`) in `http_server.go` to serve it.
- Link the stylesheet from `home.html`'s `<head>` (fragments inherit it since they're swapped into the same page).
- Apply basic layout improvements: consistent spacing/typography, a simple responsive container, clear visual separation between the National/Regional/County sections, and basic styling for tables, selects, and buttons.

## 4. Other issues worth addressing

**Resolved.** See `docs/agent/HANDOFF.md` "UI design pass".

**No loading or error feedback:** HTMX swaps show nothing while a request to met.ie is in flight, and there's no visible error message in the UI if a fetch fails — errors are handled server-side (`weather_tool/errors.go`) but never surfaced distinctly to the user. Fix: add an `hx-indicator` for in-flight requests and render a visible error fragment on failure instead of an empty/stale swap. → Added a spinner via `hx-indicator` on each trigger, and the existing `htmx:responseError` handler now renders into a styled `[role="alert"]`.

**No indication of current selection:** After choosing a county or province, nothing in the layout persists which option is currently selected, making it easy to lose track of context when comparing regions. Fix: highlight or label the active selection near the result, e.g. by echoing the selected county/province name above the forecast data. → Added a `.selected-label` per section, updated client-side on click/change.
