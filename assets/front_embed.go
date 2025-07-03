package assets

import "embed"

//go:embed all:front/build
var FrontStaticFS embed.FS
