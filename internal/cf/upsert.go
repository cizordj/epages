package cf

import (
	"context"
	"epages/internal/auth"
	"epages/internal/manifest"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func UpsertAssetHashes(
	client *cloudflare.Client,
	manifest *manifest.Manifest,
	token *auth.Token,
) error {
	_, err := client.Pages.Assets.UpsertHashes(
		context.TODO(),
		pages.AssetUpsertHashesParams{
			Hashes: cloudflare.F(
				manifest.Hashes(),
			),
		},
		option.WithHeader(
			"Authorization",
			fmt.Sprintf("Bearer %s", token.Raw),
		),
	)
	return err
}
