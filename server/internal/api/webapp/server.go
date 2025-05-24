package webapp

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

// https://vishnubharathi.codes/blog/baking-a-react-app-into-a-go-server/#Basic-serve

//go:embed dist
var WebAssets embed.FS

func New(mux *http.ServeMux) error {
	reactApp, err := fs.Sub(WebAssets, "dist")
	if err != nil {
		return fmt.Errorf("failed to find the dist directory: %w", err)
	}

	if _, err := reactApp.Open("index.html"); err != nil {
		return fmt.Errorf("index.html not found in react webapp: %w", err)
	}

	mux.Handle("/", http.FileServerFS(reactApp))
	return nil
}
