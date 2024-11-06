package data

import (
	"context"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/ONSdigital/log.go/v2/log"
	"time"
)

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
	resourceList := make([]models.Resource, numResources)
	resourceList, err = r.populateResourceList(resourceList)
	if err != nil {
		log.Error(ctx, "failed to populate resources list", err, logData)
		return nil, err
	}

	resources := &models.Resources{
		Count:        len(resourceList),
		ResourceList: resourceList,
		Limit:        option.Limit,
		Offset:       option.Offset,
		TotalCount:   numResources,
	}

	return resources, nil
}

// getResourcesCount returns the total number of jobs stored in the jobs collection in mongo
func (r *ResourceStore) getResourcesCount() (int, error) {
	return 2, nil
}

func (r *ResourceStore) populateResourceList(resourceList []models.Resource) ([]models.Resource, error) {

	topics1 := []string{"a", "b", "c", "d"}
	topics2 := []string{"a", "b", "e", "f"}
	dateChanges := []string{"a change_notice", "a previous_date"}
	tempResource1 := models.Resource{
		Uri:             "/a/temp/uri",
		UriOld:          "/an/old/uri",
		ContentType:     "api_dataset_landing_page",
		CDID:            "ASELECTIONOFNUMBERSANDLETTERS123",
		DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS456",
		Edition:         "a temporary edition",
		MetaDescription: "a temporary description",
		ReleaseDate:     time.Time{}.UTC(),
		Summary:         "a temporary summary",
		Title:           "a temporary title",
		Topics:          topics1,
		Language:        "string",
		Survey:          "string",
		CanonicalTopic:  "string",
		Cancelled:       true,
		Finalised:       true,
		Published:       true,
		DateChanges:     dateChanges,
		ProvisionalDate: "October-November 2024",
	}

	resourceList[0] = tempResource1

	tempResource2 := models.Resource{
		Uri:             "/another/temp/uri",
		UriOld:          "/another/old/uri",
		ContentType:     "api_dataset_landing_page",
		CDID:            "ASELECTIONOFNUMBERSANDLETTERS789",
		DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS101112",
		Edition:         "another temporary edition",
		MetaDescription: "another temporary description",
		ReleaseDate:     time.Time{}.UTC(),
		Summary:         "another temporary summary",
		Title:           "another temporary title",
		Topics:          topics2,
		Language:        "string",
		Survey:          "string",
		CanonicalTopic:  "string",
		Cancelled:       true,
		Finalised:       true,
		Published:       true,
		DateChanges:     dateChanges,
		ProvisionalDate: "October-November 2024",
	}

	resourceList[1] = tempResource2

	return resourceList, nil
}
