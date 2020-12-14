package cache

import (
	"context"
	"encoding/json"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const _configBucket = "configs"

type dataService struct {
	dataService avcontrol.DataService
	db          *bolt.DB
	log         *zap.Logger
}

func New(ds avcontrol.DataService, path string) (*dataService, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open cache: %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(_configBucket))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache: %v", err)
	}

	return &dataService{
		dataService: ds,
		db:          db,
	}, nil
}

func (d *dataService) RoomConfig(ctx context.Context, id string) (avcontrol.RoomConfig, error) {
	config, err := d.dataService.RoomConfig(ctx, id)
	if err != nil {
		config, cacheErr := d.roomConfigFromCache(ctx, id)
		if cacheErr != nil {
			return avcontrol.RoomConfig{}, fmt.Errorf("unable to get config from cache: %v", cacheErr)
		}

		return config, nil
	}

	if err := d.cacheConfig(ctx, id, config); err != nil {
		d.log.Warn("unable to cache config", zap.Error(err))
	}

	return config, nil
}

func (d *dataService) roomConfigFromCache(ctx context.Context, id string) (avcontrol.RoomConfig, error) {
	var config avcontrol.RoomConfig

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_configBucket))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}

		bytes := b.Get([]byte(id))
		if bytes == nil {
			return fmt.Errorf("config not in cache")
		}

		if err := json.Unmarshal(bytes, &config); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return avcontrol.RoomConfig{}, err
	}

	return config, nil
}

func (d *dataService) cacheConfig(ctx context.Context, id string, config avcontrol.RoomConfig) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_configBucket))
		if b == nil {
			return fmt.Errorf("config bucket does not exist")
		}

		bytes, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("unable to marshal config: %v", err)
		}

		if err = b.Put([]byte(id), bytes); err != nil {
			return fmt.Errorf("unable to put config: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
