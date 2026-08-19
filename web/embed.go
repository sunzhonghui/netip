package web

import "embed"

// DistFS contains the built static frontend assets.
//
//go:embed all:dist
var DistFS embed.FS
