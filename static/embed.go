package static

import "embed"

//go:embed *
var staticFiles embed.FS

func GetFiles() embed.FS {
	return staticFiles
}