package contenttype

import (
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		relPath  string
		contents []byte
		expected string
	}{
		{
			name:     "No extension, default content",
			relPath:  "file",
			contents: []byte("test content"),
			expected: "text/plain",
		},
		{
			name:     "Valid JSON content",
			relPath:  "data.json",
			contents: []byte(`{"key": "value"}`),
			expected: "application/json",
		},
		{
			name:     "No path extension, relying on content (text/html)",
			relPath:  "index",
			contents: []byte("<!DOCTYPE><html><body>Test</body></html>"),
			expected: "text/html",
		},
		{
			name:     "Empty path and empty content",
			relPath:  "",
			contents: []byte{},
			expected: "application/octet-stream",
		},
		{
			name:     "Path with dots but no recognized extension",
			relPath:  "config.foo",
			contents: []byte{},
			expected: "application/octet-stream",
		},
		{
			name:     "Path with dots but no recognized extension",
			relPath:  "image.unknown",
			contents: []byte{},
			expected: "application/octet-stream",
		},
		{
			name:     "Known extension CSS",
			relPath:  "styles.css",
			contents: []byte{},
			expected: "text/css",
		},
		{
			name:     "Known extension JS",
			relPath:  "script.js",
			contents: []byte{},
			expected: "text/javascript",
		},
		{
			name:     "Known extension MJS",
			relPath:  "module.mjs",
			contents: []byte{},
			expected: "text/javascript",
		},
		{
			name:     "Known extension JSON",
			relPath:  "data.json",
			contents: []byte{},
			expected: "application/json",
		},
		{
			name:     "Known extension HTML",
			relPath:  "index.html",
			contents: []byte{},
			expected: "text/html",
		},
		{
			name:     "Known extension HTM",
			relPath:  "index.htm",
			contents: []byte{},
			expected: "text/html",
		},
		{
			name:     "Known extension SVG",
			relPath:  "graphic.svg",
			contents: []byte{},
			expected: "image/svg+xml",
		},
		{
			name:     "Known extension XML",
			relPath:  "data.xml",
			contents: []byte{},
			expected: "text/xml",
		},
		{
			name:     "Known extension TXT",
			relPath:  "readme.txt",
			contents: []byte{},
			expected: "text/plain",
		},
		{
			name:     "Known extension WEBP",
			relPath:  "image.webp",
			contents: []byte{},
			expected: "image/webp",
		},
		{
			name:     "Known extension WOFF",
			relPath:  "font.woff",
			contents: []byte{},
			expected: "font/woff",
		},
		{
			name:     "Known extension WOFF2",
			relPath:  "font.woff2",
			contents: []byte{},
			expected: "font/woff2",
		},
		{
			name:     "Known extension TTF",
			relPath:  "font.ttf",
			contents: []byte{},
			expected: "font/ttf",
		},
		{
			name:     "Known extension MAP",
			relPath:  "bundle.map",
			contents: []byte{},
			expected: "application/json",
		},
		{
			name:     "Known extension WASM",
			relPath:  "module.wasm",
			contents: []byte{},
			expected: "application/wasm",
		},
		{
			name:     "Known extension ICO",
			relPath:  "favicon.ico",
			contents: []byte{},
			expected: "image/vnd.microsoft.icon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Detect(tt.relPath, tt.contents)
			if actual != tt.expected {
				t.Errorf(
					"Detect(%q, %v) = %q; want %q",
					tt.relPath,
					tt.contents,
					actual,
					tt.expected,
				)
			}
		})
	}
}

func TestStripParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "With semicolon parameter",
			input:    "text/html; charset=utf-8",
			expected: "text/html",
		},
		{
			name:     "Without semicolon",
			input:    "image/png",
			expected: "image/png",
		},
		{
			name:     "Multiple parameters",
			input:    "application/json; charset=utf-8; profile=v1",
			expected: "application/json",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := stripParams(tt.input)
			if actual != tt.expected {
				t.Errorf(
					"stripParams(%q) = %q; want %q",
					tt.input,
					actual,
					tt.expected,
				)
			}
		})
	}
}
