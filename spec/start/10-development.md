
## 19. Development environment

### Development mode

Run two development servers:

```text
Svelte development server: http://127.0.0.1:5173
Go API server:             http://127.0.0.1:8080
```

The Svelte development server should proxy `/api` requests to the Go server.

Example workflow:

```bash
# Terminal 1
cd frontend
npm run dev

# Terminal 2
cd backend
go run ./cmd/statescore
```

### Production mode

1. Build the Svelte static frontend.
2. Copy or generate the frontend output inside the Go embedding directory.
3. Compile the Go executable.
4. Run the executable.
5. Serve the embedded frontend and API from one localhost origin.

Example:

```bash
npm run build
go build -o statescore ./cmd/statescore
./statescore
```

---

## 20. Suggested repository structure

```text
statescore/
â”œâ”€â”€ cmd/
â”‚   â””â”€â”€ statescore/
â”‚       â””â”€â”€ main.go
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ api/
â”‚   â”œâ”€â”€ app/
â”‚   â”œâ”€â”€ browser/
â”‚   â”œâ”€â”€ config/
â”‚   â”œâ”€â”€ database/
â”‚   â”œâ”€â”€ importer/
â”‚   â”œâ”€â”€ metrics/
â”‚   â”œâ”€â”€ scoring/
â”‚   â”œâ”€â”€ security/
â”‚   â””â”€â”€ shutdown/
â”œâ”€â”€ migrations/
â”œâ”€â”€ frontend/
â”‚   â”œâ”€â”€ src/
â”‚   â”œâ”€â”€ static/
â”‚   â”œâ”€â”€ package.json
â”‚   â””â”€â”€ svelte.config.js
â”œâ”€â”€ web/
â”‚   â””â”€â”€ dist/
â”œâ”€â”€ datasets/
â”‚   â””â”€â”€ starter/
â”œâ”€â”€ scripts/
â”œâ”€â”€ tests/
â”œâ”€â”€ go.mod
â”œâ”€â”€ Makefile
â””â”€â”€ README.md
```
