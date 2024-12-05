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
func (r *ResourceStore) GetResources(ctx context.Context, options Options) (*models.Resources, error) {
	logData := log.Data{"options": options}
	log.Info(ctx, "getting list of resources", logData)

	items, err := populateItems()
	if err != nil {
		logData["items"] = items
		logData["count"] = len(items)
		log.Error(ctx, "failed to populate resources list", err, logData)
		return nil, err
	}

	filteredItems := filterItems(items, options)

	resources := &models.Resources{
		Count:      len(filteredItems),
		Items:      filteredItems,
		Limit:      options.Limit,
		Offset:     options.Offset,
		TotalCount: len(items),
	}

	log.Info(ctx, "retrieved resources", log.Data{
		"Count": len(items),
	})
	return resources, nil
}

func populateItems() (items []models.Resource, err error) {
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

// filterItems filters a list of resources by limit and offset, capping the maximum value
// at the returned items length
func filterItems(items []models.Resource, options Options) []models.Resource {
	var maxItem int
	maxRequested := options.Offset + options.Limit

	if len(items) < maxRequested {
		maxItem = len(items)
	} else {
		maxItem = maxRequested
	}

	return items[options.Offset:maxItem]
}
