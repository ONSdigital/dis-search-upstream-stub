package models

import (
	"time"
)

// Resource represents a resource metadata model and json representation for API
type Resource struct {
	URI             string    `json:"uri"`
	URIOld          string    `json:"uri_old"`
	ContentType     string    `json:"content_type"`
	CDID            string    `json:"cdid"`
	DatasetID       string    `json:"dataset_id"`
	Edition         string    `json:"edition"`
	MetaDescription string    `json:"meta_description"`
	ReleaseDate     time.Time `json:"release_date"`
	Summary         string    `json:"summary"`
	Title           string    `json:"title"`
	Topics          []string  `json:"topics"`
	Language        string    `json:"language"`
	Survey          string    `json:"survey"`
	CanonicalTopic  string    `json:"canonical_topic"`
	Cancelled       bool      `json:"cancelled"`
	Finalised       bool      `json:"finalised"`
	Published       bool      `json:"published"`
	DateChanges     []string  `json:"date_changes"`
	ProvisionalDate string    `json:"provisional_date"`
}