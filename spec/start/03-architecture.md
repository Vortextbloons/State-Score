## 5. Runtime architecture

```text
Default browser
      |
      | HTTP requests
      v
http://127.0.0.1:<port>
      |
      v
Go application
  â”œâ”€â”€ Static Svelte frontend
  â”œâ”€â”€ REST API
  â”œâ”€â”€ Scoring engine
  â”œâ”€â”€ Dataset importer
  â””â”€â”€ SQLite database
```

The frontend and API must use the same origin:

```text
Frontend: http://127.0.0.1:8787/
API:      http://127.0.0.1:8787/api/v1/
```

This avoids unnecessary cross-origin configuration.

---

## 6. Application startup behavior

When the executable starts, it must:

1. Determine the operating-system-specific application data directory.
2. Create the application directory when it does not exist.
3. Open or create the SQLite database.
4. Run pending database migrations.
5. Load configuration.
6. Attempt to bind to `127.0.0.1:8787`.
7. Use another available local port if port 8787 is occupied.
8. Start the HTTP server.
9. Wait until the server is ready.
10. Open the application URL in the default browser.
11. Continue running until the user stops the process or selects **Shut Down Application**.

The server must bind to `127.0.0.1`, not `0.0.0.0`, unless the user deliberately enables network access in a future advanced setting.

Example startup output:

```text
StateScore is running.

Open: http://127.0.0.1:8787
Data: /home/user/.local/share/statescore/statescore.db

Press Ctrl+C to stop.
```

If the browser cannot be opened automatically, the terminal must display a clickable or copyable local address.

Starting the executable a second time should either:

* Open the existing application URL, or
* Detect that another StateScore process is already using the configured port and exit cleanly.

---

## 7. Shutdown behavior

The application must support:

* `Ctrl+C` shutdown from the terminal.
* Operating-system termination signals.
* A **Shut Down Application** button in the settings page.
* Graceful completion or cancellation of active imports.
* Closing the SQLite connection.
* Graceful HTTP server shutdown.

Closing the browser tab does not automatically stop the Go process.

The interface should clearly communicate this:

```text
Closing this browser tab does not shut down StateScore.
Use Settings â†’ Shut Down Application or stop the terminal process.
```

---
