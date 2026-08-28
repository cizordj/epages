package cf

import (
	"context"
	"encoding/base64"
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
	uploadToken string,
) {
	logging.Info("uploading assets", "count", manifest.Len())

	batch := make([]pages.AssetUploadParamsBody, 0, cfg.MaxUploadCount)

	for _, entry := range manifest.All() {
		fullPath := fmt.Sprintf("%s/%s", cfg.Folder, entry.Path)
		fileBytes, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", fullPath, err)
			continue
		}

		encodedVal := base64.StdEncoding.EncodeToString(fileBytes)

		batch = append(batch, pages.AssetUploadParamsBody{
			Base64: cloudflare.F(true),
			Key:    cloudflare.F(entry.Hash),
			Value:  cloudflare.F(encodedVal),
		})

		if len(batch) >= int(cfg.MaxUploadCount) {
			sendBatch(client, batch, uploadToken)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		sendBatch(client, batch, uploadToken)
	}
	logging.Info("batch uploaded")
}

func sendBatch(client *cloudflare.Client, batch []pages.AssetUploadParamsBody, uploadToken string) {
	_, err := client.Pages.Assets.Upload(
		context.TODO(),
		pages.AssetUploadParams{
			Body: batch,
		},
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", uploadToken)),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to upload asset batch: %s", err.Error()))
	}
	logging.Debug("uploaded a batch", "count", len(batch))
}
