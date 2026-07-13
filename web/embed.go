package web

import "embed"

// Dist holds the SvelteKit static build output (see frontend adapter pages/assets).
//
//go:embed all:dist
var Dist embed.FS
