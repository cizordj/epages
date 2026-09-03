package cf

import (
	"context"
	"encoding/json"
	"epages/internal/config"
	"epages/internal/gitinfo"
	"epages/internal/logging"
	"epages/internal/manifest"
	"os"
	"path/filepath"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

type DeploymentInfo struct {
	Url     string
	Id      string
	Aliases []string
}

func CreateDeployment(
	client *cloudflare.Client,
	cfg *config.Config,
	m *manifest.Manifest,
	excludedFiles *manifest.Manifest,
) (*DeploymentInfo, error) {
	jsonBytes, err := json.Marshal(
		m.Exclude(excludedFiles.Hashes()).APIManifest(),
	)
	if err != nil {
		return &DeploymentInfo{}, err
	}

	params := pages.ProjectDeploymentNewParams{
		AccountID: cloudflare.F(cfg.AccountId),
		Manifest:  cloudflare.F(string(jsonBytes)),
	}

	headersEntry, ok := m.ByPath("_headers")

	if ok {
		relPath := filepath.Join(cfg.Folder, headersEntry.Path)
		fileHandle, err := os.Open(relPath)
		if err != nil {
			logging.Warn(err)
		} else {
			params.Headers = cloudflare.FileParam(
				fileHandle,
				filepath.Base(headersEntry.Path),
				"text/plain",
			)
		}
	}

	redirectsEntry, ok := m.ByPath("_redirects")

	if ok {
		relPath := filepath.Join(cfg.Folder, redirectsEntry.Path)
		fileHandle, err := os.Open(relPath)
		if err != nil {
			logging.Warn(err)
		} else {
			params.Redirects = cloudflare.FileParam(
				fileHandle,
				filepath.Base(redirectsEntry.Path),
				"text/plain",
			)
		}
	}

	pwd, err := os.Getwd()

	if err != nil {
		logging.Warn(err.Error())
	} else {
		repoInfo, err := gitinfo.Discover(pwd)

		if err != nil {
			logging.Warn(
				"could not discover git repository information",
				"error",
				err,
			)
		} else {
			params.Branch = cloudflare.F(
				repoInfo.Branch,
			)
			params.CommitHash = cloudflare.F(
				repoInfo.Hash,
			)
			params.CommitMessage = cloudflare.F(
				repoInfo.Message,
			)
			if repoInfo.Dirty {
				params.CommitDirty = cloudflare.F(
					pages.ProjectDeploymentNewParamsCommitDirtyTrue,
				)
			} else {
				params.CommitDirty = cloudflare.F(
					pages.ProjectDeploymentNewParamsCommitDirtyFalse,
				)
			}
		}
	}

	deployment, err := client.Pages.Projects.Deployments.New(
		context.TODO(),
		cfg.ProjectName,
		params,
	)
	if err != nil {
		return &DeploymentInfo{}, err
	}

	return &DeploymentInfo{
		Url:     deployment.URL,
		Id:      deployment.ID,
		Aliases: deployment.Aliases,
	}, nil
}
