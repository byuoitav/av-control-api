package init

import (
	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/common/log"
)

//DefaultInitializer implements the Initializer interface
type DefaultInitializer struct {
}

//Initialize fulfills the initializers for the Initializer interface
func (i *DefaultInitializer) Initialize(room base.Room) error {
	log.L.Info("[init] Yay! I work.\n")
	return nil
}
