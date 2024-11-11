package models

// Resources represents an array of Resource objects and json representation for API
type Resources struct {
	Count        int           `json:"count"`
	ResourceList []interface{} `json:"resources"`
	Limit        int           `json:"limit"`
	Offset       int           `json:"offset"`
	TotalCount   int           `json:"total_count"`
}
