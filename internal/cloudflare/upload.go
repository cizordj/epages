package upload

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func GetUploadToken(client *cloudflare.Client, projectName *string, accountId *string) string {
	response, err := client.Pages.Projects.GetUploadToken(
		context.TODO(),
		*projectName,
		pages.ProjectGetUploadTokenParams{
			AccountID: cloudflare.F(*accountId),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	return response.JWT
}

func CheckMissingAssets(client *cloudflare.Client, hashes []string, uploadToken string) []string {
	page, err := client.Pages.Assets.CheckMissing(context.TODO(), pages.AssetCheckMissingParams{
		Hashes: cloudflare.F(hashes),
	},
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", uploadToken)),
	)
	if err != nil {
		panic(err.Error())
	}
	return page.Result
}

func UploadAssets(client *cloudflare.Client, hashFileMap map[string]string, uploadToken string, rootDir string) {
	const batchSize = 50
	
	batch := make([]pages.AssetUploadParamsBody, 0, batchSize)

	for hash, relPath := range hashFileMap {
		fullPath := fmt.Sprintf("%s/%s", rootDir, relPath)
		fileBytes, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", fullPath, err)
			continue
		}

		encodedVal := base64.StdEncoding.EncodeToString(fileBytes)

		batch = append(batch, pages.AssetUploadParamsBody{
			Base64: cloudflare.F(true),
			Key:    cloudflare.F(hash),
			Value:  cloudflare.F(encodedVal),
		})

		if len(batch) >= batchSize {
			sendBatch(client, batch, uploadToken)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		sendBatch(client, batch, uploadToken)
	}
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
	fmt.Printf("Successfully uploaded a batch of %d assets.\n", len(batch))
}
