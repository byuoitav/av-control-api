package couch

import "testing"

func TestOptions(t *testing.T) {
	opts := []Option{
		WithBasicAuth("user", "pass"),
		WithInsecure(),
		WithDatabase("db"),
		WithMappingDocumentID("docID"),
		WithEnvironment("env"),
	}

	options := options{
		scheme:       _defaultScheme,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
	}

	for _, o := range opts {
		o.apply(&options)
	}
}
