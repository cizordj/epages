package cf

import (
	"context"
	"encoding/json"
	"epage/internal/config"
	"epage/internal/manifest"

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
	manifest *manifest.Manifest,
) (*DeploymentInfo, error) {
	jsonBytes, err := json.Marshal(manifest.APIManifest())
	if err != nil {
		return &DeploymentInfo{}, err
	}

	deployment, err := client.Pages.Projects.Deployments.New(
		context.TODO(),
		cfg.ProjectName,
		pages.ProjectDeploymentNewParams{
			AccountID: cloudflare.F(cfg.AccountId),
			Manifest:  cloudflare.F(string(jsonBytes)),
		},
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
