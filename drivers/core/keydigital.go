package core

type KeyDigitalDriver struct{}

func (k *KeyDigitalDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

// func (k *KeyDigitalDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
// 	return keydigital.CreateVideoSwitcher(ctx, addr, drivers.Log.Named(addr))
// }
