package upload

import (
	"context"
	"fmt"

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

func UploadAssets(client *cloudflare.Client, hashFileMap map[string]string, uploadToken string) {
	response, err := client.Pages.Assets.Upload(context.TODO(), pages.AssetUploadParams{
		Body: []pages.AssetUploadParamsBody{{
			Base64: cloudflare.F(true),
			Key:    cloudflare.F("b026324c6904b2a9cb4b88d6d61c81d1"),
			Metadata: cloudflare.F(pages.AssetUploadParamsBodyMetadata{
				ContentType: cloudflare.F("text/plain"),
			}),
			Value: cloudflare.F("SGVsbG8sIFdvcmxkIQ=="),
		}},
	},
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", uploadToken)),
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", response.Errors)
}
