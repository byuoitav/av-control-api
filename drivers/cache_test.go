package drivers

import (
	"context"
	"math/rand"
	"testing"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
)

type testDriver struct {
	createDelay func() time.Duration
}

func (d *testDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (d *testDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	time.Sleep(d.createDelay())
	return &struct {
		// makes sure that we don't return the same address every time
		// https://golang.org/ref/spec#Size_and_alignment_guarantees
		// > zero-size variables may have the same address
		field string
	}{}, nil
}

func TestSavingDevices(t *testing.T) {
	is := is.New(t)
	rand.Seed(time.Now().Unix())

	cache := &deviceCache{
		Driver: &testDriver{
			createDelay: func() time.Duration {
				return time.Duration(rand.Intn(50)) * time.Millisecond
			},
		},
		cache: make(map[string]avcontrol.Device),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := "1.1.1.1"
	dev, err := cache.CreateDevice(ctx, addr)
	is.NoErr(err)

	for i := 0; i < 100; i++ {
		d, err := cache.CreateDevice(ctx, addr)
		is.NoErr(err)
		is.True(d == dev)
	}

	// try a few random addresses, make sure they are different
	d, err := cache.CreateDevice(ctx, "1.0.0.1")
	is.NoErr(err)
	is.True(d != dev)

	d, err = cache.CreateDevice(ctx, "0.0.0.0")
	is.NoErr(err)
	is.True(d != dev)

	d, err = cache.CreateDevice(ctx, "random string")
	is.NoErr(err)
	is.True(d != dev)
}

func TestSaveDevicesAtSameTime(t *testing.T) {
	is := is.New(t)

	cache := &deviceCache{
		Driver: &testDriver{
			createDelay: func() time.Duration {
				return 2 * time.Second
			},
		},
		cache: make(map[string]avcontrol.Device),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	group, gctx := errgroup.WithContext(ctx)
	addr := "1.1.1.1"

	devs := make([]avcontrol.Device, 3)
	done := make([]time.Time, 3)

	group.Go(func() error {
		var err error
		devs[0], err = cache.CreateDevice(gctx, addr)
		done[0] = time.Now()
		return err
	})

	group.Go(func() error {
		time.Sleep(50 * time.Millisecond)

		var err error
		devs[1], err = cache.CreateDevice(gctx, addr)
		done[1] = time.Now()
		return err
	})

	group.Go(func() error {
		time.Sleep(100 * time.Millisecond)

		var err error
		devs[2], err = cache.CreateDevice(gctx, addr)
		done[2] = time.Now()
		return err
	})

	is.NoErr(group.Wait())

	// make sure all the devs are the same
	for i := 0; i < len(devs)-2; i++ {
		is.True(devs[i] == devs[i+1])
	}

	// make sure they showed up at about the same time
	for i := 0; i < len(done)-2; i++ {
		cur := done[i].Round(25 * time.Millisecond)
		next := done[i].Round(25 * time.Millisecond)

		is.Equal(cur, next)
	}
}
