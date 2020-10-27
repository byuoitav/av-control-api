package drivers

/*
// TODO this should test the deviceCache
func TestSavingDevices(t *testing.T) {
	rand.Seed(time.Now().Unix())

	newDev := saveDevicesFunc(func(ctx context.Context, addr string) (Device, error) {
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		return &mock.TV{}, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := "1.1.1.1"
	dev, err := newDev(ctx, addr)
	if err != nil {
		t.Fatalf("unable to create device: %s", err)
	}

	for i := 0; i < 20; i++ {
		d, err := newDev(ctx, addr)
		if err != nil {
			t.Fatalf("unable to create device: %s", err)
		}

		if d != dev {
			t.Fatalf("got mismatched devices: expected %p, got %p", dev, d)
		}
	}

	// try a few random addresses, make sure they are different
	d, err := newDev(ctx, "")
	if err != nil {
		t.Fatalf("unable to create device: %s", err)
	}

	if d == dev {
		t.Fatalf("got the same devices for different addresses: got %p", d)
	}

	// try with 0.0.0.0
	d, err = newDev(ctx, "0.0.0.0")
	if err != nil {
		t.Fatalf("unable to create device: %s", err)
	}

	if d == dev {
		t.Fatalf("got the same devices for different addresses: got %p", d)
	}

	// try with one more
	d, err = newDev(ctx, "this is a string!")
	if err != nil {
		t.Fatalf("unable to create device: %s", err)
	}

	if d == dev {
		t.Fatalf("got the same devices for different addresses: got %p", d)
	}
}

func TestSaveDevicesAtSameTime(t *testing.T) {
	newDev := saveDevicesFunc(func(ctx context.Context, addr string) (Device, error) {
		time.Sleep(2 * time.Second)
		return &mock.TV{}, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	group, gctx := errgroup.WithContext(ctx)
	addr := "1.1.1.1"

	devs := make([]Device, 3)
	done := make([]time.Time, 3)

	group.Go(func() error {
		var err error
		devs[0], err = newDev(gctx, addr)
		if err != nil {
			return err
		}

		done[0] = time.Now()
		return nil
	})

	group.Go(func() error {
		time.Sleep(50 * time.Millisecond)

		var err error
		devs[1], err = newDev(gctx, addr)
		if err != nil {
			return err
		}

		done[1] = time.Now()
		return nil
	})

	group.Go(func() error {
		time.Sleep(100 * time.Millisecond)

		var err error
		devs[2], err = newDev(gctx, addr)
		if err != nil {
			return err
		}

		done[2] = time.Now()
		return nil
	})

	err := group.Wait()
	if err != nil {
		t.Fatalf("unexpected error creating devices: %s", err)
	}

	// make sure all the devs are the same
	if devs[0] != devs[1] || devs[0] != devs[2] {
		t.Fatalf("not all devices matched. got %p, %p, and %p", devs[0], devs[1], devs[2])
	}

	// make sure they all showed up at about the same time
	rounded := make([]time.Time, 3)
	rounded[0] = done[0].Round(25 * time.Millisecond)
	rounded[1] = done[1].Round(25 * time.Millisecond)
	rounded[2] = done[2].Round(25 * time.Millisecond)

	if !rounded[0].Equal(rounded[1]) || !rounded[0].Equal(rounded[2]) {
		t.Fatalf("didn't finish around the same time. finished at %v, %v, and %v", done[0].Format(time.RFC3339Nano), done[1].Format(time.RFC3339Nano), done[2].Format(time.RFC3339Nano))
	}
}
*/
