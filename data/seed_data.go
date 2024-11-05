package data

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Highlight represents the structure of the highlight field.
type Highlight struct {
	Summary string `json:"summary"`
	Title   string `json:"title"`
}

// ContentItem represents a dataset item structure with various metadata fields.
type ContentItem struct {
	Type            string     `json:"type"`
	DatasetID       string     `json:"dataset_id"`
	Keywords        []string   `json:"keywords,omitempty"`
	MetaDescription string     `json:"meta_description"`
	ReleaseDate     *time.Time `json:"release_date,omitempty"` // Use pointer for nullable dates
	Summary         string     `json:"summary"`
	Title           string     `json:"title"`
	URI             string     `json:"uri"`
	Highlight       Highlight  `json:"highlight"` // Nested struct for highlights
	Topics          []string   `json:"topics"`
	CanonicalTopic  string     `json:"canonical_topic"`
}

//go:embed *.json
var dataFiles embed.FS

func main() {
	var contentItems []ContentItem

	// Get a list of JSON files in the embedded 'data' directory.
	files, err := dataFiles.ReadDir("data")
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
	}

	// Loop through files, read, and unmarshal each JSON file into Go structs.
	for _, file := range files {
		f, err := dataFiles.ReadFile("data/" + file.Name())
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", file.Name(), err)
		}

		var fileResources []ContentItem
		// Unmarshal JSON arrays into slices
		if err := json.Unmarshal(f, &fileResources); err != nil {
			log.Fatalf("Failed to unmarshal JSON for file %s: %v", file.Name(), err)
		}
		contentItems = append(contentItems, fileResources...)
	}

	fmt.Printf("Loaded %d content items\n", len(contentItems))

	// Display the loaded data with readable formatting
	for _, item := range contentItems {
		// Pretty-print each item as JSON for verification
		data, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			log.Fatalf("Failed to format content item: %v", err)
		}
		fmt.Println(string(data))
	}
}
