package main

import (
	"epages/internal/cf"
	"epages/internal/config"
	"epages/internal/logging"
	"epages/internal/manifest"
)

func init() {
	config.DefineFlags()
}

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		logging.Fatal(err)
	}

	files, err := manifest.GenerateHashMap(cfg.Folder)

	if err != nil {
		logging.Fatal(err)
	}

	ignoredFiles := manifest.GetIgnoredFiles(files)

	logging.Info("detected files for upload", "count", files.Len()-ignoredFiles.Len())

	client := cf.WithClient(cfg)

	uploadToken, err := cf.GenUploadToken(client, cfg)

	if err != nil {
		logging.Fatal(err)
	}

	missingFiles, err := cf.CheckMissingAssets(
		client,
		uploadToken,
		files,
		ignoredFiles,
	)

	if err != nil {
		logging.Fatal(err)
	}

	cf.UploadAssets(client, missingFiles, cfg, uploadToken)

	deploymentInfo, err := cf.CreateDeployment(client, cfg, files, ignoredFiles)

	if err != nil {
		logging.Fatal(err)
	}

	uploadToken, err = cf.RefreshTokenIfExpired(client, cfg, uploadToken)

	if err != nil {
		logging.Fatal(err)
	}

	err = cf.UpsertAssetHashes(client, files, uploadToken)

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
