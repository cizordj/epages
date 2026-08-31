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
