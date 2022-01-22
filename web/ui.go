package web

import "embed"

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

//go:embed build
var FS embed.FS
