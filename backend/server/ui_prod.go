//go:build serveui

package server

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed dist/*
var uiFS embed.FS

func GetUiHandler() (http.Handler, error) {
	dist, err := fs.Sub(uiFS, "dist")
	if err != nil {
		return nil, fmt.Errorf("Failed to get dist subdirectory: %w", err)
	}

	return &uiHandler{fs: dist}, nil
}

type uiHandler struct{ fs fs.FS }

func (h *uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if path == "" {
		h.serveIndex(w, r)
		return
	}

	file, err := h.fs.Open(path)
	if err != nil {
		h.serveIndex(w, r)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		h.serveIndex(w, r)
		return
	}

	if info.IsDir() {
		h.serveIndex(w, r)
		return
	}

	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file.(io.ReadSeeker))
}

func (h *uiHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	indexContent, err := fs.ReadFile(h.fs, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(indexContent)
}
