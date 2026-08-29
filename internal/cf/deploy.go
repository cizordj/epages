package cf

import (
	"context"
	"encoding/json"
	"epage/internal/config"
	"epage/internal/logging"
	"epage/internal/manifest"
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
			params.Headers = cloudflare.FileParam(
				fileHandle,
				filepath.Base(redirectsEntry.Path),
				"text/plain",
			)
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
