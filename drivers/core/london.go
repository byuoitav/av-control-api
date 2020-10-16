package core

type LondonDriver struct{}

func (l *LondonDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

// func (l *LondonDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
// 	return london.New(addr, london.WithLogger(drivers.Log.Named(addr))), nil
// }
