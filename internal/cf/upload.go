package cf

import (
	"context"
	"encoding/base64"
	"epage/internal/auth"
	"epage/internal/config"
	"epage/internal/logging"
	"epage/internal/manifest"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func UploadAssets(
	client *cloudflare.Client,
	manifest *manifest.Manifest,
	cfg *config.Config,
	token *auth.Token,
) {
	logging.Info("start uploading assets")

	batch := make([]pages.AssetUploadParamsBody, 0, cfg.MaxUploadCount)

	for _, entry := range manifest.All() {
		fullPath := fmt.Sprintf("%s/%s", cfg.Folder, entry.Path)
		fileBytes, err := os.ReadFile(fullPath)
		if err != nil {
			logging.Warn(
				"Warning: failed to read file",
				"file",
				fullPath,
				"error",
				err.Error(),
			)
			continue
		}

		encodedVal := base64.StdEncoding.EncodeToString(fileBytes)

		batch = append(batch, pages.AssetUploadParamsBody{
			Base64: cloudflare.F(true),
			Key:    cloudflare.F(entry.Hash),
			Value:  cloudflare.F(encodedVal),
		})

		if len(batch) >= int(cfg.MaxUploadCount) {
			sendBatch(client, batch, token)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		sendBatch(client, batch, token)
	}
}

func sendBatch(client *cloudflare.Client, batch []pages.AssetUploadParamsBody, token *auth.Token) {
	_, err := client.Pages.Assets.Upload(
		context.TODO(),
		pages.AssetUploadParams{
			Body: batch,
		},
		option.WithHeader(
			"Authorization",
			fmt.Sprintf("Bearer %s", token.Raw),
		),
	)
	if err != nil {
		logging.Error(
			"Failed to upload asset batch",
			"error",
			err.Error(),
		)
	}
	logging.Debug("uploaded a batch", "count", len(batch))
}
