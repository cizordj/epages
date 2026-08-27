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

func CheckMissingAssets(client *cloudflare.Client, hashes []string, uploadToken string) {
	page, err := client.Pages.Assets.CheckMissing(context.TODO(), pages.AssetCheckMissingParams{
		Hashes: cloudflare.F(hashes),
	},
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", uploadToken)),
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", page)
}
