package cf

import (
	"context"
	"epage/internal/config"
	"epage/internal/logging"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func GetUploadToken(client *cloudflare.Client, cfg *config.Config) (string, error) {
	response, err := client.Pages.Projects.GetUploadToken(
		context.TODO(),
		cfg.ProjectName,
		pages.ProjectGetUploadTokenParams{
			AccountID: cloudflare.F(cfg.AccountId),
		},
	)
	if err != nil {
		return "", err
	}
	logging.Debug("upload token", "jwt", response.JWT)
	return response.JWT, nil
}
