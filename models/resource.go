package models

// Resource represents a standard resource metadata model and json representation for API
type Resource struct {
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
	Title           string   `avro:"title" json:"title"`
	Topics          []string `avro:"topics" json:"topics"`
	URI             string   `avro:"uri" json:"uri"`
	URIOld          string   `avro:"uri_old" json:"uri_old"`
	Release         Release  `avro:"release" json:"release"`
}

// Release contains the additional resource fields that are only used for content_type=release
type Release struct {
	Cancelled       bool     `avro:"cancelled,omitempty" json:"cancelled,omitempty"`
	Finalised       bool     `avro:"finalised,omitempty" json:"finalised,omitempty"`
	Published       bool     `avro:"published,omitempty" json:"published,omitempty"`
	DateChanges     []string `avro:"date_changes,omitempty" json:"date_changes,omitempty"`
	ProvisionalDate string   `avro:"provisional_date,omitempty" json:"provisional_date,omitempty"`
}
