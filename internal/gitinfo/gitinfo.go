package gitinfo

import (
	"os/exec"
	"strings"
)

type Info struct {
	Branch  string
	Hash    string
	Message string
	Dirty   bool
}

// Discover returns best-effort git metadata for dir. If dir isn't a git
// repo (or git isn't installed), it returns a zero Info and a nil error —
// git metadata is optional, not fatal.
func Discover(dir string) Info {
	branch, _ := run(dir, "rev-parse", "--abbrev-ref", "HEAD")
	hash, _ := run(dir, "rev-parse", "HEAD")
	msg, _ := run(dir, "log", "-1", "--pretty=%s")
	status, _ := run(dir, "status", "--porcelain")
	return Info{
		Branch:  branch,
		Hash:    hash,
		Message: msg,
		Dirty:   status != "",
	}
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
