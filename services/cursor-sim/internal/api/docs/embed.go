package docs

import "embed"

// staticFiles contains embedded static assets including OpenAPI specs
//
//go:embed static/*
var staticFiles embed.FS
