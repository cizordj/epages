package contenttype

import (
	"os/exec"
	"testing"
)

func TestOSFileStrategy(t *testing.T) {
	strategy := &OSFileStdinStrategy{}
	_, err := exec.LookPath("file")
	if err != nil {
		t.Skip(err)
	}

	tests := []struct {
		name     string
		relPath  string
		contents []byte
		expected string
	}{
		{
			name:     "Makefile",
			relPath:  "Makefile",
			contents: []byte("all:\n\techo 123\nclean:\n\trm -rf bin"),
			expected: "text/x-makefile",
		},
		{
			name:     "GPG public key",
			relPath:  "public_key.asc",
			contents: []byte("-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n"),
			expected: "application/pgp-keys",
		},
		{
			name:     "File with empty contents",
			relPath:  "empty.txt",
			contents: []byte(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := strategy.Detect(
				tt.relPath,
				tt.contents,
			)

			if actual != tt.expected {
				t.Errorf(
					"OSFileStdinStrategy.Detect(%q, %v) = %q; want %q",
					tt.relPath,
					tt.contents,
					actual,
					tt.expected,
				)
			}
		})
	}
}

func TestSkipOSFileStrategy(t *testing.T) {
	t.Setenv("PATH", "")
	t.Run(
		"",
		func(t *testing.T) {
			strategy := &OSFileStdinStrategy{}
			_, ok := strategy.Detect("foo.txt", []byte("hello world"))
			if ok != false {
				t.Error(
					"os file strategy should be ignored when the `file` command is not present",
				)
			}
		},
	)
}
