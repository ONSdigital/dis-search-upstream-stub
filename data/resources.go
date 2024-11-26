package data

import (
	"context"
	"embed"

	"encoding/json"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

//go:embed json_files/*.json
var jsonFiles embed.FS

// GetResources retrieves all the resources from the collection
func (r *ResourceStore) GetResources(ctx context.Context, option Options) (*models.Resources, error) {
	logData := log.Data{"option": option}
	log.Info(ctx, "getting list of resources", logData)

	items, err := r.populateItems()
	if err != nil {
		logData["items"] = items
		logData["count"] = len(items)
		log.Error(ctx, "failed to populate resources list", err, logData)
		return nil, err
	}

	resources := &models.Resources{
		Count:      len(items),
		Items:      items,
		Limit:      option.Limit,
		Offset:     option.Offset,
		TotalCount: len(items),
	}

	log.Info(ctx, "retrieved resoucres", log.Data{
		"Count": len(items),
	})
	return resources, nil
}

func (r *ResourceStore) populateItems() (items []models.Resource, err error) {
	// Get a list of JSON files and/or subdirectories in the embedded 'json_files' directory.
	dirEntries, err := jsonFiles.ReadDir("json_files")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read json_files directory")
	}
	items = make([]models.Resource, 0, len(dirEntries))

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		fileBytes, err := jsonFiles.ReadFile("json_files/" + dirEntry.Name())
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file")
		}

		var resource models.Resource
		err = json.Unmarshal(fileBytes, &resource)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal JSON for file")
		}

		items = append(items, resource)
	}

	return items, nil
}
