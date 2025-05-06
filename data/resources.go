package data

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"encoding/json"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:embed json_files/content_updated/*.json json_files/search_content_updated/*.json
var jsonFiles embed.FS

const (
	searchContentUpdatedResourceType = "SearchContentUpdatedResource"
	contentUpdatedResourceType       = "ContentUpdatedResource"
)

// GetResources is the method that satisfies the DataStorer interface
// It calls the existing GetResourcesWithType with a default resourceType
func (r *ResourceStore) GetResources(ctx context.Context, options Options) (*models.Resources, error) {
	// Use a default resourceType (or it could be dynamic based on your use case)
	resourceType := "searchContentUpdatedResourceType" // Set this based on your needs or pass an empty string or default value

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

	// Set the directory based on resource type
	if resourceType == contentUpdatedResourceType {
		dir = "json_files/content_updated"
	} else if resourceType == searchContentUpdatedResourceType {
		dir = "json_files/search_content_updated"
	} else {
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	// Read files from the appropriate directory
	dirEntries, err := fs.ReadDir(jsonFiles, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	items := make([]models.Resource, 0, len(dirEntries))

	// Loop through files in the directory
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue // Skip directories, only process files
		}

		// Read the content of each file
		fileBytes, err := jsonFiles.ReadFile(dir + "/" + dirEntry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", dirEntry.Name(), err)
		}

		var resource models.Resource
		// Unmarshal the file content based on the requested resource type
		if resourceType == contentUpdatedResourceType {
			var contentUpdatedResource models.ContentUpdatedResource
			err = json.Unmarshal(fileBytes, &contentUpdatedResource)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal ContentUpdatedResource for file %s: %w", dirEntry.Name(), err)
			}
			resource = contentUpdatedResource
		} else if resourceType == searchContentUpdatedResourceType {
			var searchContentUpdatedResource models.SearchContentUpdatedResource
			err = json.Unmarshal(fileBytes, &searchContentUpdatedResource)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal SearchContentUpdatedResource for file %s: %w", dirEntry.Name(), err)
			}
			resource = searchContentUpdatedResource
		}

		// Append the resource to the items slice
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
