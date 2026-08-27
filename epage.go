package main

import (
	"epage/internal/cloudflare/deployment"
	"epage/internal/cloudflare/upload"
	"epage/internal/hasher"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/spf13/pflag"
)

var (
	token       *string
	projectName *string
	accountId   *string
	folder      *string
)

func init() {
	token = pflag.StringP("api-token", "t", os.Getenv("CLOUDFLARE_API_TOKEN"), "The Cloudflare API Token (can also be set via CLOUDFLARE_API_TOKEN env var)")
	projectName = pflag.StringP("project-name", "n", os.Getenv("CLOUDFLARE_PROJECT_NAME"), "The name of the Cloudflare Pages project (can also be set via CLOUDFLARE_PROJECT_NAME env var)")
	accountId = pflag.StringP("account-id", "a", os.Getenv("CLOUDFLARE_ACCOUNT_ID"), "The Cloudflare account ID (can also be set via CLOUDFLARE_ACCOUNT_ID env var)")
	folder = pflag.StringP("folder", "f", "", "The folder containing the assets to deploy")
}

func main() {
	pflag.Parse()

	if len(os.Args) == 1 || (len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		pflag.Usage()
		return
	}

	var missingFlags []string
	if *token == "" {
		missingFlags = append(missingFlags, "api-token")
	}
	if *projectName == "" {
		missingFlags = append(missingFlags, "project-name")
	}
	if *accountId == "" {
		missingFlags = append(missingFlags, "account-id")
	}
	if *folder == "" {
		missingFlags = append(missingFlags, "folder")
	}

	if len(missingFlags) > 0 {
		fmt.Printf("Error: Missing required flags: --%s\n", strings.Join(missingFlags, ", --"))
		fmt.Println("\nUse -h or --help for usage information.")
		os.Exit(1)
	}

	missingFlags = nil

	var _, folderErr = os.ReadDir(*folder)

	if folderErr != nil {
		fmt.Printf("%s\n", folderErr.Error())
		os.Exit(1)
	}

	fileHashMap, err := hasher.GenerateFileHashMap(*folder)
	if err != nil {
		fmt.Printf("Error generating file hashes: %s\n", err.Error())
		os.Exit(1)
	}

	client := cloudflare.NewClient(
		option.WithAPIToken(*token),
	)

	uploadToken := upload.GetUploadToken(client, projectName, accountId)

	setOfHashes := make([]string, 0, len(fileHashMap))
	for hash := range fileHashMap {
		setOfHashes = append(setOfHashes, hash)
	}

	missingAssetsHash := upload.CheckMissingAssets(client, setOfHashes, uploadToken)
	newFileHashMap := make(map[string]string)

	for _, hash := range missingAssetsHash {
		newFileHashMap[hash] = fileHashMap[hash]
	}

	upload.UploadAssets(client, newFileHashMap, uploadToken, *folder)
	deployment.CreateDeployment(client, projectName, accountId, fileHashMap)
}
