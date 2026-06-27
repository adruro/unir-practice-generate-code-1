package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed all:templates
var TemplatesFS embed.FS

//go:embed all:static
var StaticFS embed.FS

func GetStaticSubFS() fs.FS {
	sub, err := fs.Sub(StaticFS, "static")
	if err != nil {
		log.Fatalf("Error getting static sub FS: %v", err)
	}
	return sub
}
