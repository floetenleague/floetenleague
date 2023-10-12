package ui

import (
	"embed"
	"io/fs"
)

//go:embed build
var buildFiles embed.FS
var Files, _ = fs.Sub(buildFiles, "build")
