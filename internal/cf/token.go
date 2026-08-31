package cf

import (
	"context"
	"epages/internal/auth"
	"epages/internal/config"
	"epages/internal/logging"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func GenUploadToken(client *cloudflare.Client, cfg *config.Config) (*auth.Token, error) {
	response, err := client.Pages.Projects.GetUploadToken(
		context.TODO(),
		cfg.ProjectName,
		pages.ProjectGetUploadTokenParams{
			AccountID: cloudflare.F(cfg.AccountId),
		},
	)
	if err != nil {
		return nil, err
	}
	token, err := auth.New(response.JWT)
	if err != nil {
		return nil, err
	}
	logging.Info("successfully generated the upload token")
	return token, nil
}

func RefreshTokenIfExpired(client *cloudflare.Client, cfg *config.Config, auth *auth.Token) (*auth.Token, error) {
	if auth.IsExpired() {
		logging.Debug("upload token is expired, getting a new one")
		return GenUploadToken(
			client,
			cfg,
		)
	} else {
		logging.Debug("upload token is not expired, reusing the same one")
		return auth, nil
	}
}
