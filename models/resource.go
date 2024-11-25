package models

import (
	"time"
)

// Standard represents a standard resource metadata model and json representation for API
type Resource struct {
	CanonicalTopic  string    `json:"canonical_topic"`
	CDID            string    `json:"cdid"`
	ContentType     string    `json:"content_type"`
	DatasetID       string    `json:"dataset_id"`
	Edition         string    `json:"edition"`
	Language        string    `json:"language"`
	MetaDescription string    `json:"meta_description"`
	ReleaseDate     time.Time `json:"release_date"`
	Summary         string    `json:"summary"`
	Survey          string    `json:"survey"`
	Title           string    `json:"title"`
	Topics          []string  `json:"topics"`
	URI             string    `json:"uri"`
	URIOld          string    `json:"uri_old"`
	Release
}

// These fields are only used for content_type=release
type Release struct {
	Cancelled       bool     `json:"cancelled,omitempty"`
	Finalised       bool     `json:"finalised,omitempty"`
	Published       bool     `json:"published,omitempty"`
	DateChanges     []string `json:"date_changes,omitempty"`
	ProvisionalDate string   `json:"provisional_date,omitempty"`
}
