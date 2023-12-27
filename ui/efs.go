package ui

import (
	"embed"
)

// commments that start with go: is a special comment directive

//go:embed "html" "static"
var Files embed.FS

// this embeds static resources (CSS, JS, SQL files) into the binary
