package data

import "embed"

//go:embed names.json distros.dhall
var FS embed.FS
