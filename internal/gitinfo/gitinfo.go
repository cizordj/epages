package gitinfo

import (
	"strings"

	"github.com/go-git/go-git/v6"
)

type Info struct {
	Branch  string
	Hash    string
	Message string
	Dirty   bool
}

func Discover(dir string) (Info, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return Info{}, err
	}

	head, err := repo.Head()
	if err != nil {
		return Info{}, err
	}

	var branchName string
	if head.Name().IsBranch() {
		branchName = head.Name().Short()
	} else {
		branchName = "HEAD (Detached: " + head.Hash().String()[:7] + ")"
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return Info{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return Info{}, err
	}
	status, err := worktree.Status()
	if err != nil {
		return Info{}, err
	}

	return Info{
		Branch:  branchName,
		Hash:    commit.Hash.String(),
		Message: strings.TrimSpace(commit.Message),
		Dirty:   !status.IsClean(),
	}, nil
}
