package cache

import (
	"context"
	"fmt"
	"os"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
)

type mockDataService struct {
	configs map[string]avcontrol.RoomConfig
}

func (m *mockDataService) RoomConfig(ctx context.Context, id string) (avcontrol.RoomConfig, error) {
	config, ok := m.configs[id]
	if !ok {
		return config, fmt.Errorf("config not found")
	}

	return config, nil
}

func TestCache(t *testing.T) {

	testConfig := avcontrol.RoomConfig{
		ID: "test",
		Devices: map[avcontrol.DeviceID]avcontrol.DeviceConfig{
			"device": avcontrol.DeviceConfig{
				Address: "hello.com",
				Driver:  "adam",
			},
		},
	}

	mock := &mockDataService{
		configs: map[string]avcontrol.RoomConfig{
			"test": testConfig,
		},
	}

	file := os.TempDir() + "/av-control-api-cache-test.db"
	ds, err := New(mock, file)
	require.NoError(t, err)
	defer os.Remove(file)

	is := is.New(t)

	t.Run("ConfigPassThrough", func(t *testing.T) {
		config, err := ds.RoomConfig(context.TODO(), "test")
		is.NoErr(err)
		is.Equal(config, testConfig)
	})

	t.Run("ConfigCached", func(t *testing.T) {
		delete(mock.configs, "test")

		config, err := ds.RoomConfig(context.TODO(), "test")
		is.NoErr(err)
		is.Equal(config, testConfig)
	})

	t.Run("ConfigMissing", func(t *testing.T) {
		_, err := ds.RoomConfig(context.TODO(), "config")
		is.True(err != nil)
	})
}
