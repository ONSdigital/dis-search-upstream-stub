package data

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"encoding/json"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

//go:embed json_files/search_content_updated/*.json json_files/search_content_deleted/*.json json_files/content_updated/*.json
var jsonFiles embed.FS

var searchContentUpdatedResourceType = "SearchContentUpdatedResource"
var searchContentDeletedResourceType = "SearchContentDeletedResource"
var contentUpdatedResourceType = "ContentUpdatedResource"

// GetResources is the method that satisfies the DataStorer interface
// It calls the existing GetResourcesWithType with a default resourceType
func (r *ResourceStore) GetResources(ctx context.Context, typeParam string, options Options) (*models.Resources, error) {
	// Use a default resourceType (or it could be dynamic based on your use case)
	var resourceType string

	switch typeParam {
	case "content-updated":
		resourceType = contentUpdatedResourceType
	case "search-content-deleted":
		resourceType = searchContentDeletedResourceType
	default:
		resourceType = searchContentUpdatedResourceType
	}

	// Call the existing method with resourceType
	return r.GetResourcesWithType(ctx, resourceType, options)
}

// GetResourcesWithType retrieves all the resources from the collection
func (r *ResourceStore) GetResourcesWithType(ctx context.Context, resourceType string, options Options) (*models.Resources, error) {
	logData := log.Data{"options": options}
	log.Info(ctx, "getting list of resources", logData)

	items, err := populateItems(resourceType)
	if err != nil {
		logData["items"] = items
		logData["count"] = len(items)
		logData["type"] = resourceType
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

// populateItems retrieves items from the content_updated and search_content_updated directories
func populateItems(resourceType string) ([]models.Resource, error) {
	var dir string

	// Determine which directory to read from based on the resource type
	switch resourceType {
	case contentUpdatedResourceType:
		dir = "json_files/content_updated"
	case searchContentUpdatedResourceType:
		dir = "json_files/search_content_updated"
	case searchContentDeletedResourceType:
		dir = "json_files/search_content_deleted"
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	// Read files from the appropriate directory
	dirEntries, err := fs.ReadDir(jsonFiles, dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read json_files directory")
	}

	items := make([]models.Resource, 0, len(dirEntries))

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		// Read the content of each file
		fileBytes, err := jsonFiles.ReadFile(dir + "/" + dirEntry.Name())
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file")
		}

		var resource models.Resource

		// Determine which type of resource to unmarshal into
		switch resourceType {
		case contentUpdatedResourceType:
			var contentUpdated models.ContentUpdatedResource
			err = json.Unmarshal(fileBytes, &contentUpdated)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal ContentUpdatedResource JSON")
			}
			resource = contentUpdated
		case searchContentUpdatedResourceType:
			var searchContentUpdated models.SearchContentUpdatedResource
			err = json.Unmarshal(fileBytes, &searchContentUpdated)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal SearchContentUpdatedResource JSON")
			}
			resource = searchContentUpdated
		case searchContentDeletedResourceType:
			var searchContentDeleted models.SearchContentDeletedResource
			err = json.Unmarshal(fileBytes, &searchContentDeleted)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal SearchContentDeletedResource JSON")
			}
			resource = searchContentDeleted
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
