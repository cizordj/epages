package config

import (
	"epages/internal/logging"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	token          *string
	projectName    *string
	accountId      *string
	folder         *string
	logLevel       *string
	maxUploadCount *uint16
	versionFlag    *bool
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type Config struct {
	Token          string
	ProjectName    string
	AccountId      string
	Folder         string
	MaxUploadCount uint16
}

func DefineFlags() {
	token = pflag.StringP(
		"api-token",
		"t",
		os.Getenv("CLOUDFLARE_API_TOKEN"),
		"The Cloudflare API Token (can also be set via CLOUDFLARE_API_TOKEN env var)",
	)
	projectName = pflag.StringP(
		"project-name",
		"n",
		os.Getenv("CLOUDFLARE_PROJECT_NAME"),
		"The name of the Cloudflare Pages project (can also be set via CLOUDFLARE_PROJECT_NAME env var)",
	)
	accountId = pflag.StringP(
		"account-id",
		"a",
		os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		"The Cloudflare account ID (can also be set via CLOUDFLARE_ACCOUNT_ID env var)",
	)
	logLevel = pflag.StringP(
		"log-level",
		"l",
		"off",
		"Levels, least to most verbose: off, error, warn, info, debug, trace",
	)
	folder = pflag.StringP(
		"folder",
		"f",
		"",
		"The folder containing the assets to deploy",
	)
	maxUploadCount = pflag.Uint16P(
		"max-upload-count",
		"m",
		uint16(500),
		"The maximum number of files to be uploaded at once",
	)
	versionFlag = pflag.BoolP(
		"version",
		"v",
		false,
		"Prints the version",
	)
}

func ParseConfig() (*Config, error) {
	pflag.Parse()

	if len(os.Args) == 1 || (len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help")) {
		pflag.Usage()
		return &Config{}, errors.New("")
	}

	if *versionFlag {
		fmt.Printf("epage %s (commit %s, built %s)\n", version, commit, date)
		os.Exit(0)
	}

	logLevel, err := logging.ParseLevel(*logLevel)

	if err != nil {
		return &Config{}, err
	}
	logging.Init(logLevel)

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
		return &Config{}, fmt.Errorf(
			"Missing required flags: --%s\n%s\n",
			strings.Join(missingFlags, ", --"),
			"Use -h or --help for usage information.",
		)
	}

	missingFlags = nil

	logging.Info(
		"starting epages",
		"version",
		version,
		"commit",
		commit,
		"date",
		date,
	)

	return &Config{
		Token:          *token,
		ProjectName:    *projectName,
		AccountId:      *accountId,
		Folder:         *folder,
		MaxUploadCount: *maxUploadCount,
	}, nil
}
