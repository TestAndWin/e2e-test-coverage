package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var assets embed.FS

// AssetFile returns the file system for the embedded assets
// It strips the "dist" prefix so that the router sees the files at the root
func AssetFile() http.FileSystem {
	fsys, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

// MustAsset returns the content of the index.html file
func MustAsset(name string) []byte {
	// fs.Sub treats the path relative to "dist"
	// So we can read directly from the stripped fs
	fsys, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}

	bytes, err := fs.ReadFile(fsys, name)
	if err != nil {
		panic(err)
	}
	return bytes
}
