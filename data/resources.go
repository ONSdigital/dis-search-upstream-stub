package data

import (
	"context"
	"embed"
	"fmt"

	"encoding/json"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:embed json_files/*.json
var jsonFiles embed.FS

// GetResources retrieves all the resources from the collection
func (r *ResourceStore) GetResources(ctx context.Context, option Options) (*models.Resources, error) {
	logData := log.Data{"option": option}
	log.Info(ctx, "getting list of resources", logData)

	// get resources count
	numResources, err := r.getResourcesCount()
	if err != nil {
		log.Error(ctx, "failed to get resources count", err, logData)
		return nil, err
	}

	// create and populate resourcesList
	items := make([]interface{}, numResources)
	items, err = r.populateItems(ctx, items, numResources)
	if err != nil {
		log.Error(ctx, "failed to populate resources list", err, logData)
		return nil, err
	}

	resources := &models.Resources{
		Count:      len(items),
		Items:      items,
		Limit:      option.Limit,
		Offset:     option.Offset,
		TotalCount: numResources,
	}

	return resources, nil
}

// getResourcesCount returns the total number of resources stored
func (r *ResourceStore) getResourcesCount() (int, error) {
	return 3, nil
}

func (r *ResourceStore) populateItems(ctx context.Context, items []interface{}, numItems int) ([]interface{}, error) {
	// Get a list of JSON files in the embedded 'json_files' directory.
	files, err := jsonFiles.ReadDir("json_files")
	if err != nil {
		log.Fatal(ctx, "Failed to read json_files directory", err)
	}

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for i := 0; i < numItems; i++ {
		file := files[i]
		fmt.Println("Now reading file: " + file.Name())
		fileBytes, err := jsonFiles.ReadFile("json_files/" + file.Name())
		if err != nil {
			log.Fatal(ctx, "Failed to read file: "+file.Name(), err)
		}

		var jsonData map[string]interface{}
		err = json.Unmarshal(fileBytes, &jsonData)
		if err != nil {
			log.Fatal(ctx, "Failed to unmarshal JSON for file: "+file.Name(), err)
		}

		items[i] = jsonData
	}

	return items, nil
}
