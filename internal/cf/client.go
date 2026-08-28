package cf

import (
	"epage/internal/config"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

func WithClient(cfg *config.Config) *cloudflare.Client {
	return cloudflare.NewClient(
		option.WithAPIToken(cfg.Token),
	)
}
