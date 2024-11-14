package data

// ResourceStore is a type that contains an implementation of the DataStorer interface, which can be used for
// getting Resources.
type ResourceStore struct {
}

// Options contains information for pagination which includes offset and limit
type Options struct {
	Offset int
	Limit  int
}
