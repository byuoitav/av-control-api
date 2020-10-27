package cache

import (
	"context"
	"fmt"
	"os"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
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
		ID: "yourmom",
		Devices: map[avcontrol.DeviceID]avcontrol.DeviceConfig{
			"your mom": avcontrol.DeviceConfig{
				Address: "hello.com",
				Driver:  "adam",
			},
		},
	}

	mock := &mockDataService{
		configs: map[string]avcontrol.RoomConfig{
			"yourmom": testConfig,
		},
	}

	file := os.TempDir() + "/av-control-api-cache-test.db"
	ds, err := New(mock, file)
	require.NoError(t, err)
	defer os.Remove(file)

	t.Run("ConfigPassThrough", func(t *testing.T) {
		config, err := ds.RoomConfig(context.TODO(), "yourmom")
		require.NoError(t, err)
		require.Equal(t, config, testConfig)
	})

	t.Run("ConfigCached", func(t *testing.T) {
		delete(mock.configs, "yourmom")

		config, err := ds.RoomConfig(context.TODO(), "yourmom")
		require.NoError(t, err)
		require.Equal(t, config, testConfig)
	})

	t.Run("ConfigMissing", func(t *testing.T) {
		_, err := ds.RoomConfig(context.TODO(), "mymom")
		require.Error(t, err)
	})
}
