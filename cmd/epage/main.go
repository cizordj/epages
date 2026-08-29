package main

import (
	"epage/internal/cf"
	"epage/internal/config"
	"epage/internal/logging"
	"epage/internal/manifest"
)

func init() {
	config.DefineFlags()
}

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		logging.Fatal(err)
	}

	manifest, err := manifest.GenerateHashMap(cfg.Folder)

	if err != nil {
		logging.Fatal(err)
	}
	logging.Info("detected files for upload", "count", manifest.Len())

	client := cf.WithClient(cfg)

	uploadToken, err := cf.GenUploadToken(client, cfg)

	if err != nil {
		logging.Fatal(err)
	}

	missingFiles, err := cf.CheckMissingAssets(client, uploadToken, manifest)

	if err != nil {
		logging.Fatal(err)
	}

	cf.UploadAssets(client, missingFiles, cfg, uploadToken)

	deploymentInfo, err := cf.CreateDeployment(client, cfg, manifest)

	if err != nil {
		logging.Fatal(err)
	}

	uploadToken, err = cf.RefreshTokenIfExpired(client, cfg, uploadToken)

	if err != nil {
		logging.Fatal(err)
	}

	err = cf.UpsertAssetHashes(client, missingFiles, uploadToken)

	if err != nil {
		logging.Fatal(err)
	}

	logging.Info(
		"succesfully deployed to cloudflare pages",
		"url",
		deploymentInfo.Url,
		"id",
		deploymentInfo.Id,
		"aliases",
		deploymentInfo.Aliases,
	)
}
