package contenttype

import (
	"bytes"
	"context"
	"epages/internal/logging"
	"mime"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Strategy interface {
	Detect(relPath string, contents []byte) (string, bool)
	Name() string
}

func (s *OSFileStdinStrategy) Name() string { return "OSFileStdinStrategy" }
func (s *MimePackageStrategy) Name() string { return "MimePackageStrategy" }
func (s *HttpDetectStrategy) Name() string  { return "HttpDetectStrategy" }

type MimePackageStrategy struct{}

func (s *MimePackageStrategy) Detect(relPath string, contents []byte) (string, bool) {
	ext := filepath.Ext(relPath)
	if ext == "" {
		return "", false
	}
	if t := mime.TypeByExtension(ext); t != "" {
		return stripParams(t), true
	}
	return "", false
}

type HttpDetectStrategy struct{}

func (s *HttpDetectStrategy) Detect(relPath string, contents []byte) (string, bool) {
	if len(contents) == 0 {
		return "", false
	}
	t := http.DetectContentType(contents)
	mimeType := stripParams(t)
	if mimeType == "" {
		return "", false
	}
	return mimeType, true
}

type OSFileStdinStrategy struct{}

func (s *OSFileStdinStrategy) Detect(relPath string, contents []byte) (string, bool) {
	if len(contents) == 0 {
		logging.Debug("OSFileStdinStrategy: empty contents, skipping")
		return "", false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "file", "--mime-type", "-b", "-")
	cmd.Stdin = bytes.NewReader(contents)

	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logging.Debug("OSFileStdinStrategy: `file` command timed out")
		} else {
			logging.Debug("OSFileStdinStrategy: `file` command failed", "error", err)
		}
		return "", false
	}

	mimeType := strings.TrimSpace(string(out))
	if mimeType == "" || mimeType == "application/octet-stream" {
		logging.Debug("OSFileStdinStrategy: treating as no-match", "mimetype", mimeType)
		return "", false
	}
	return mimeType, true
}

var weakMimeTypes = map[string]bool{
	"application/octet-stream": true,
	"text/plain":               true,
}

func init() {
	mime.AddExtensionType(".map", "application/json")
	if path, err := exec.LookPath("file"); err == nil {
		logging.Debug(
			"`file` binary found",
			"path",
			path,
		)
	} else {
		logging.Debug("`file` binary not found in PATH, OSFileStdinStrategy will be skipped for all inputs")
	}
}

func Detect(relPath string, contents []byte) string {
	var fallback string
	strategies := []Strategy{
		&OSFileStdinStrategy{},
		&MimePackageStrategy{},
		&HttpDetectStrategy{},
	}

	logging.Debug("Detect:", "relPath", relPath, "contentsLen", len(contents))

	for _, strategy := range strategies {
		mimeType, ok := strategy.Detect(relPath, contents)

		if !ok {
			logging.Debug("Strategy did not match", "name", strategy.Name())
			continue
		}

		if weakMimeTypes[mimeType] {
			logging.Debug(
				"weak match, continuing to next strategy",
				"name",
				strategy.Name(),
				"mimetype",
				mimeType,
			)
			if fallback == "" {
				fallback = mimeType
			}
			continue
		}

		logging.Debug("Matched", "name", strategy.Name(), "mimetype", mimeType)
		return mimeType
	}

	if fallback != "" {
		logging.Debug("Detect: no strong match, using fallback", "mimetype", fallback)
		return fallback
	}

	logging.Debug(
		"Detect: no strategy matched, using defaults",
		"mimetype",
		"application/octet-stream",
	)
	return "application/octet-stream"
}

func stripParams(t string) string {
	before, _, _ := strings.Cut(t, ";")
	return strings.TrimSpace(before)
}
