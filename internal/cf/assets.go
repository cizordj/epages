package cf

import (
	"context"
	"epage/internal/auth"
	"epage/internal/logging"
	"epage/internal/manifest"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func CheckMissingAssets(client *cloudflare.Client, uploadToken *auth.Token, man *manifest.Manifest) (*manifest.Manifest, error) {
	page, err := client.Pages.Assets.CheckMissing(
		context.TODO(),
		pages.AssetCheckMissingParams{
			Hashes: cloudflare.F(man.Hashes()),
		},
		option.WithHeader(
			"Authorization",
			fmt.Sprintf("Bearer %s", uploadToken.Raw),
		),
	)
	if err != nil {
		return manifest.New(), err
	}
	logging.Info(
		"files that are missing on the edge",
		"count",
		len(page.Result),
	)
	return man.Subset(page.Result), nil
}
