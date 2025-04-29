package models

// Resource interface represents a standard resource metadata model and json representation for API
type Resource interface {
	GetResourceType() string
}
