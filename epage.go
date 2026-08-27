package main

import (
	"epage/internal/hasher"
	"fmt"
	"os"
	"strings"

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

	fileHashes, err := hasher.GenerateFileHashMap(*folder)
	if err != nil {
		fmt.Printf("Error generating file hashes: %s\n", err.Error())
		os.Exit(1)
	}

	for path, hash := range fileHashes {
		fmt.Printf("%s: %s\n", path, hash)
	}
}
