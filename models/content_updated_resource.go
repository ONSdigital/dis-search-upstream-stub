package models

// ContentUpdatedResource represents a standard resource metadata model and json representation for API
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
