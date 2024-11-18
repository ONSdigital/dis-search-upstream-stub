package models

// Resources represents an array of Standard and Release resource objects and json representation for API
type Resources struct {
	Count      int           `json:"count"`
	Items      []interface{} `json:"items"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	TotalCount int           `json:"total_count"`
}
