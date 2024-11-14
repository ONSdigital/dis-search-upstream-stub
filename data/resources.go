package data

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// //go:embed standard_resource_1.json
// var resource1 []byte
//
// //go:embed standard_resource_2.json
// var resource2 []byte
//
// //go:embed standard_resource_3.json
// var resource3 []byte

////go:embed folder/*.hash
//var folder embed.FS

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
	// items, err = r.populateResourceList(items)
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

//	func (r *ResourceStore) populateResourceList(items []interface{}) ([]interface{}, error) {
//		topics1 := []string{"a", "b", "c", "d"}
//		topics2 := []string{"a", "b", "e", "f"}
//		dateChanges := []string{"a change_notice", "a previous_date"}
//		tempResource1 := models.Standard{
//			URI:             "/a/temp/uri",
//			URIOld:          "/an/old/uri",
//			ContentType:     "api_dataset_landing_page",
//			CDID:            "ASELECTIONOFNUMBERSANDLETTERS123",
//			DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS456",
//			Edition:         "a temporary edition",
//			MetaDescription: "a temporary description",
//			ReleaseDate:     time.Time{}.UTC(),
//			Summary:         "a temporary summary",
//			Title:           "a temporary title",
//			Topics:          topics1,
//			Language:        "string",
//			Survey:          "string",
//			CanonicalTopic:  "string",
//		}
//
//		items[0] = tempResource1
//
//		tempResource2 := models.Release{
//			URI:             "/another/temp/uri",
//			URIOld:          "/another/old/uri",
//			ContentType:     "api_dataset_landing_page",
//			CDID:            "ASELECTIONOFNUMBERSANDLETTERS789",
//			DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS101112",
//			Edition:         "another temporary edition",
//			MetaDescription: "another temporary description",
//			ReleaseDate:     time.Time{}.UTC(),
//			Summary:         "another temporary summary",
//			Title:           "another temporary title",
//			Topics:          topics2,
//			Language:        "string",
//			Survey:          "string",
//			CanonicalTopic:  "string",
//			Cancelled:       true,
//			Finalised:       true,
//			Published:       true,
//			DateChanges:     dateChanges,
//			ProvisionalDate: "October-November 2024",
//		}
//
//		items[1] = tempResource2
//
//		return items, nil
//	}
func (r *ResourceStore) populateItems(ctx context.Context, items []interface{}, numItems int) ([]interface{}, error) {
	// Get a list of JSON files in the embedded 'json_files' directory.
	files, err := jsonFiles.ReadDir("json_files")
	if err != nil {
		log.Fatal(ctx, "Failed to read json_files directory", err)
	}

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for i := 0; i < numItems; i++ {
		//startIndex := numItems - 1
		//for i := startIndex; i >= 0; i-- {
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

		//sort.Sort(sort.Reverse(sort.StringSlice(jsonData)))
		items[i] = jsonData
	}

	return items, nil
}
