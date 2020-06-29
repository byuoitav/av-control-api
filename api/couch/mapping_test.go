package couch

import (
	"context"
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"github.com/go-kivik/kivikmock/v3"
	"github.com/google/go-cmp/cmp"
)

func TestMapping(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
		environment:  "default",
	}

	db := mock.NewDB()
	mock.ExpectDB().WithName(ds.database).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, `{
		"drivers": {
			"Sony Bravia": {
				"envs": {
					"default": {
						"address": "localhost:9001",
						"ssl": false
					}
				}
			},
			"Sony ADCP": {
				"envs": {
					"default": {
						"address": "localhost:9002",
						"ssl": true
					}
				}
			}
		}
	}`))

	mapping, err := ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected := api.DriverMapping{
		"Sony Bravia": api.DriverConfig{
			Address: "localhost:9001",
			SSL:     false,
		},
		"Sony ADCP": api.DriverConfig{
			Address: "localhost:9002",
			SSL:     true,
		},
	}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect mapping (-want, +got):\n%s", diff)
	}
}

func TestMappingMissingDrivers(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
		environment:  "default",
	}

	db := mock.NewDB()
	mock.ExpectDB().WithName(_defaultDatabase).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, `{
		"drivers": {
			"Sony Bravia": {
				"envs": {
					"default": {
						"address": "localhost:9001",
						"ssl": false
					}
				}
			},
			"Sony ADCP": {
				"envs": {
					"test": {
						"address": "localhost:9002",
						"ssl": true
					}
				}
			}
		}
	}`))

	mapping, err := ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected := api.DriverMapping{
		"Sony Bravia": api.DriverConfig{
			Address: "localhost:9001",
			SSL:     false,
		},
	}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect mapping (-want, +got):\n%s", diff)
	}
}

func TestMappingEmpty(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
		environment:  "default",
	}

	db := mock.NewDB()
	mock.ExpectDB().WithName(_defaultDatabase).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, `{
		"drivers": {
			"Sony Bravia": {
				"envs": {
					"test": {
						"address": "localhost:9001",
						"ssl": false
					}
				}
			},
			"Sony ADCP": {
				"envs": {
					"test": {
						"address": "localhost:9002",
						"ssl": true
					}
				}
			}
		}
	}`))

	mapping, err := ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected := api.DriverMapping{}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect mapping (-want, +got):\n%s", diff)
	}
}

func TestMappingMultipleEnvironments(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
		environment:  "default",
	}

	doc := `{
		"drivers": {
			"Sony Bravia": {
				"envs": {
					"default": {
						"address": "localhost:9001",
						"ssl": false
					},
					"test": {
						"address": "localhost:8001",
						"ssl": true
					}
				}
			},
			"Sony ADCP": {
				"envs": {
					"default": {
						"address": "localhost:9002",
						"ssl": true
					},
					"test": {
						"address": "localhost:8002",
						"ssl": false
					}
				}
			},
			"Atlona 5x1": {
				"envs": {
					"test": {
						"address": "localhost:8003",
						"ssl": true
					}
				}
			}
		}

	}`

	db := mock.NewDB()
	mock.ExpectDB().WithName(_defaultDatabase).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, doc))
	mock.ExpectDB().WithName(_defaultDatabase).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, doc))

	mapping, err := ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected := api.DriverMapping{
		"Sony Bravia": api.DriverConfig{
			Address: "localhost:9001",
			SSL:     false,
		},
		"Sony ADCP": api.DriverConfig{
			Address: "localhost:9002",
			SSL:     true,
		},
	}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect *default* mapping (-want, +got):\n%s", diff)
	}

	ds.environment = "test"

	mapping, err = ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected = api.DriverMapping{
		"Sony Bravia": api.DriverConfig{
			Address: "localhost:8001",
			SSL:     true,
		},
		"Sony ADCP": api.DriverConfig{
			Address: "localhost:8002",
			SSL:     false,
		},
		"Atlona 5x1": api.DriverConfig{
			Address: "localhost:8003",
			SSL:     true,
		},
	}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect *test* mapping (-want, +got):\n%s", diff)
	}
}

func TestMappingDifferentDocDB(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     "testDB",
		mappingDocID: "testMappingDocID",
		environment:  "default",
	}

	db := mock.NewDB()
	mock.ExpectDB().WithName(ds.database).WillReturn(db)
	db.ExpectGet().WithDocID(ds.mappingDocID).WillReturn(kivikmock.DocumentT(t, `{
		"drivers": {
			"Sony Bravia": {
				"envs": {
					"default": {
						"address": "localhost:9001",
						"ssl": false
					}
				}
			},
			"Sony ADCP": {
				"envs": {
					"default": {
						"address": "localhost:9002",
						"ssl": true
					}
				}
			}
		}
	}`))

	mapping, err := ds.DriverMapping(context.Background())
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expected := api.DriverMapping{
		"Sony Bravia": api.DriverConfig{
			Address: "localhost:9001",
			SSL:     false,
		},
		"Sony ADCP": api.DriverConfig{
			Address: "localhost:9002",
			SSL:     true,
		},
	}

	if diff := cmp.Diff(expected, mapping); diff != "" {
		t.Errorf("generated incorrect mapping (-want, +got):\n%s", diff)
	}
}
