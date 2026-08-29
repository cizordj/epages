package contenttype

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func init() {
	types := map[string]string{
		".css":   "text/css",
		".js":    "application/javascript",
		".mjs":   "application/javascript",
		".json":  "application/json",
		".html":  "text/html",
		".htm":   "text/html",
		".svg":   "image/svg+xml",
		".xml":   "application/xml",
		".txt":   "text/plain",
		".webp":  "image/webp",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".otf":   "font/otf",
		".map":   "application/json",
		".wasm":  "application/wasm",
		".ico":   "image/x-icon",
	}
	for ext, typ := range types {
		_ = mime.AddExtensionType(ext, typ)
	}
}

func Detect(relPath string, contents []byte) string {
	ext := filepath.Ext(relPath)
	if ext != "" {
		if t := mime.TypeByExtension(ext); t != "" {
			return stripParams(t)
		}
	}
	if len(contents) > 0 {
		return stripParams(http.DetectContentType(contents))
	}
	return "application/octet-stream"
}

func stripParams(t string) string {
	before, _, _ := strings.Cut(t, ";")
	return strings.TrimSpace(before)
}
