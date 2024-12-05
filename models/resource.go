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
