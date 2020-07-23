package couch

import "testing"

func TestOptions(t *testing.T) {
	opts := []Option{
		WithBasicAuth("user", "pass"),
		WithDatabase("db"),
		WithMappingDocumentID("docID"),
		WithEnvironment("env"),
	}

	options := options{
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
	}

	for _, o := range opts {
		o.apply(&options)
	}
}
