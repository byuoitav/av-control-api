package state

import (
	"context"
	"testing"
	"time"
)

func (t *testing.T) TestGetPowerSimple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var getPower getPower
}
