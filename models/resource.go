package models

import (
	"encoding/json"
	"fmt"
)

// Resource interface that both types will implement
type Resource interface {
	GetResourceType() string
}

// Custom unmarshalling for Resources
func (r *Resources) UnmarshalJSON(data []byte) error {
	// Define a temporary structure for unmarshalling
	var aux struct {
		Count      int               `json:"count"`
		TotalCount int               `json:"total_count"`
		Limit      int               `json:"limit"`
		Offset     int               `json:"offset"`
		Items      []json.RawMessage `json:"items"` // Store raw JSON for each item
	}

	// First, unmarshal the non-Items fields
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Now we need to unmarshal the Items into a slice of Resources (deciding the type dynamically)
	var items []Resource

	// Loop through each item and determine its concrete type based on some field
	for _, itemRaw := range aux.Items {
		var tempItem map[string]interface{}
		if err := json.Unmarshal(itemRaw, &tempItem); err != nil {
			return err
		}

		// Check for a distinguishing field (e.g., resourceType)
		// Determine the resource type based on fields present in the item
		if _, ok := tempItem["content_type"]; ok {
			var resource SearchContentUpdatedResource
			if err := json.Unmarshal(itemRaw, &resource); err != nil {
				return err
			}
			items = append(items, resource)
		} else if _, ok := tempItem["data_type"]; ok {
			var resource ContentUpdatedResource
			if err := json.Unmarshal(itemRaw, &resource); err != nil {
				return err
			}
			items = append(items, resource)
		} else {
			// Handle unknown resource type or invalid format
			return fmt.Errorf("invalid resource type")
		}
	}

	// Set the values in the Resources struct
	r.Count = aux.Count
	r.TotalCount = aux.TotalCount
	r.Limit = aux.Limit
	r.Offset = aux.Offset
	r.Items = items

	return nil
}

// ContentUpdatedResource represents the first type
type ContentUpdatedResource struct {
	URI          string `avro:"uri" json:"uri"`
	DataType     string `avro:"data_type" json:"data_type"`
	CollectionID string `avro:"collection_id" json:"collection_id"`
	JobID        string `avro:"job_id" json:"job_id"`
	SearchIndex  string `avro:"search_index" json:"search_index"`
	TraceID      string `avro:"trace_id" json:"trace_id"`
}

func (r ContentUpdatedResource) GetResourceType() string {
	return "ContentUpdatedResource"
}

// SearchContentUpdatedResource represents a standard resource metadata model and json representation for API
type SearchContentUpdatedResource struct {
	CanonicalTopic  string   `avro:"canonical_topic" json:"canonical_topic"`
	CDID            string   `avro:"cdid" json:"cdid"`
	ContentType     string   `avro:"content_type" json:"content_type"`
	DatasetID       string   `avro:"dataset_id" json:"dataset_id"`
	Edition         string   `avro:"edition" json:"edition"`
	Language        string   `avro:"language" json:"language"`
	MetaDescription string   `avro:"meta_description" json:"meta_description"`
	ReleaseDate     string   `avro:"release_date" json:"release_date"`
	Summary         string   `avro:"summary" json:"summary"`
	Survey          string   `avro:"survey" json:"survey"`
	SearchIndex     string   `avro:"search_index" json:"search_index"`
	TraceID         string   `avro:"trace_id" json:"trace_id"`
	Title           string   `avro:"title" json:"title"`
	Topics          []string `avro:"topics" json:"topics"`
	URI             string   `avro:"uri" json:"uri"`
	URIOld          string   `avro:"uri_old" json:"uri_old"`
	// These fields are only used for content_type=release
	Cancelled       bool                 `avro:"cancelled" json:"cancelled"`
	DateChanges     []ReleaseDateDetails `avro:"date_changes" json:"date_changes"`
	Finalised       bool                 `avro:"finalised" json:"finalised"`
	ProvisionalDate string               `avro:"provisional_date" json:"provisional_date"`
	Published       bool                 `avro:"published" json:"published"`
}

type ReleaseDateDetails struct {
	ChangeNotice string `avro:"change_notice" json:"change_notice"`
	PreviousDate string `avro:"previous_date" json:"previous_date"`
}

func (r SearchContentUpdatedResource) GetResourceType() string {
	return "SearchContentUpdatedResource"
}
