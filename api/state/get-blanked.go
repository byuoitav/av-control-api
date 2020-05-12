package state

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type getBlanked struct {
}

func (*getBlanked) GenerateActions(ctx context.Context, room []api.Device, env string) ([]action, []api.DeviceStateError, []DeviceStateUpdate) {
	var acts []action
	var errs []api.DeviceStateError
	var expectedUpdates []DeviceStateUpdate

	// just doing basic get blanked for now
	for _, dev := range room {
		url, order, err := getCommand(dev, "GetBlanked", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: err.Error(),
			})

			continue
		}

		// replace values
		params := map[string]string{
			"address": dev.Address,
		}
		url, err = fillURL(url, params)
		if err != nil {
			errs = append(errs, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		// build http request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			errs = append(errs, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: fmt.Sprintf("unable to build http request: %s", err),
			})

			continue
		}

		acts = append(acts, action{
			ID:    dev.ID,
			Req:   req,
			Order: order,
		})

		expectedUpdates = append(expectedUpdates, DeviceStateUpdate{
			ID: dev.ID,
		})
	}

	return acts, errs, expectedUpdates
}
