package ui

// This file helps to embed the VUE app in the binary

//go:generate npm run build
//go:generate go-bindata -fs -o ui_gen.go -pkg ui -prefix dist/ ./dist/...
