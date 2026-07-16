<img src="app-icon.png" alt="StateScore" width="200">

# StateScore Frontend

SvelteKit 2 static SPA (Svelte 5 runes, TypeScript) for ranking and comparing U.S. states.

Built with `@sveltejs/adapter-static`, outputs to `../web/dist/` for embedding into the Go binary.

## Develop

```sh
npm run dev
```

Runs on `http://127.0.0.1:5173`, proxies `/api` to `http://127.0.0.1:8080`.

## Build

```sh
npm run build
```

Outputs to `../web/dist/`.

## Check & Test

```sh
npm run check        # svelte-kit sync + svelte-check
npm test             # vitest
npm run lint         # prettier + eslint
```

## Rankings

Current default-profile rankings are shown in the root [README.md](../README.md).
