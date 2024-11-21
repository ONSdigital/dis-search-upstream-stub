package data

import (
	"context"
	"embed"
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

	items, numResources, err := r.populateItems(ctx)
	if err != nil {
		logData["items"] = items
		logData["numResources"] = numResources
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

func (r *ResourceStore) populateItems(ctx context.Context) (items []interface{}, numItems int, err error) {
	logData := log.Data{}
	// Get a list of JSON files in the embedded 'json_files' directory.
	files, err := jsonFiles.ReadDir("json_files")
	if err != nil {
		log.Fatal(ctx, "Failed to read json_files directory", err)
	}
	for range files {
		numItems++
	}
	items = make([]interface{}, numItems)

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for i := 0; i < numItems; i++ {
		file := files[i]
		logData["file"] = file.Name()
		fileBytes, err := jsonFiles.ReadFile("json_files/" + file.Name())
		if err != nil {
			log.Fatal(ctx, "Failed to read file", err, logData)
		}

		var jsonData map[string]interface{}
		err = json.Unmarshal(fileBytes, &jsonData)
		if err != nil {
			log.Fatal(ctx, "Failed to unmarshal JSON for file: "+file.Name(), err)
		}

		items[i] = jsonData
	}

	return items, numItems, nil
}
