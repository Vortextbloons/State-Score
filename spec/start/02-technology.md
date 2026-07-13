---

## 4. Technology stack

### Frontend

* SvelteKit
* TypeScript
* Static SvelteKit build
* Standard CSS or a small component library
* Charting library for comparison and trend charts

SvelteKit will use its static adapter so the frontend can be built into static HTML, JavaScript, and CSS files.

### Backend

* Go
* Go `net/http`
* JSON REST API
* Background data-import jobs
* Local scoring engine

Goâ€™s standard HTTP package will serve both the application frontend and the JSON API.

### Storage

* SQLite
* Local database file
* SQL migrations
* Optional SQLite WAL mode
* CSV and JSON import support

### Distribution

The production build should contain:

* One Go executable.
* Embedded frontend files.
* An automatically created SQLite database.
* Optional starter datasets.

The built Svelte files can be compiled into the Go executable with Goâ€™s `embed` package.

The project will not use Electron, Wails, Tauri, or another embedded-webview framework.

---
