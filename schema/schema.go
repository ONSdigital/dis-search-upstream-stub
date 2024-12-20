package schema

import "github.com/ONSdigital/dp-kafka/v4/avro"

var searchContentUpdate = `{
  "type": "record",
  "name": "search-content-updated",
  "fields": [
    {"name": "canonical_topic", "type": "string", "default": ""},
    {"name": "cdid", "type": "string", "default": ""},
    {"name": "content_type", "type": "string", "default": ""},
    {"name": "dataset_id", "type": "string", "default": ""},
    {"name": "edition", "type": "string", "default": ""},
    {"name": "language", "type": "string", "default": ""},
    {"name": "meta_description", "type": "string", "default": ""},
    {"name": "release_date", "type": "string", "default": ""},
    {"name": "summary", "type": "string", "default": ""},
    {"name": "survey", "type": "string", "default": ""},
    {"name": "title", "type": "string", "default": ""},
    {"name": "topics", "type": {"type": "array", "items": "string"}, "default": []},
    {"name": "uri", "type": "string", "default": ""},
    {"name": "uri_old", "type": "string", "default": ""},
    {"name": "cancelled", "type": "boolean", "default": false},
    {"name": "finalised", "type": "boolean", "default": false},
    {"name": "published", "type": "boolean", "default": false},
    {"name": "date_changes", "type": { "type": "array", "items": {
      "name": "ReleaseDateDetails",
      "type": "record",
      "fields": [
        { "name":"change_notice", "type":"string" },
        { "name":"previous_date", "type":"string" }
      ] 
    }}},
    {"name": "provisional_date", "type": "string", "default": ""}
  ]
}`

// SearchContentUpdateEvent the Avro schema for search-content-updated messages.
var SearchContentUpdateEvent = &avro.Schema{
	Definition: searchContentUpdate,
}
