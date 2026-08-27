package deployment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/pages"
)

func CreateDeployment(
	client *cloudflare.Client,
	projectName *string,
	accountId *string,
	fileHashMap map[string]string,
) {
	invertedMap := make(map[string]string, len(fileHashMap))

	for hash, path := range fileHashMap {
		invertedMap[path] = hash
	}
	jsonBytes, err := json.Marshal(invertedMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize manifest to JSON: %s", err.Error()))
	}

	deployment, err := client.Pages.Projects.Deployments.New(
		context.TODO(),
		*projectName,
		pages.ProjectDeploymentNewParams{
			AccountID: cloudflare.F(*accountId),
			Manifest:  cloudflare.F(string(jsonBytes)),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deployment successful: %+v\n", deployment.ID)
}
