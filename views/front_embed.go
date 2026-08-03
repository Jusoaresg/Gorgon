package views

import "embed"

//go:embed all:static
var FrontStaticFS embed.FS

//go:embed all:ui
var FrontFS embed.FS
