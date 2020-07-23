package state

import (
	"sort"

	"github.com/byuoitav/av-control-api/api"
)

// containsString checks if the given slice of strings contains the provided string.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// sortErrors sorts a slice of api.DeviceStateError. Sorting is done in this order:
//   1. ID
//   2. Field
//   3. Error
func sortErrors(errors []api.DeviceStateError) {
	sort.Slice(errors, func(i, j int) bool {
		if errors[i].ID != errors[j].ID {
			return errors[i].ID < errors[j].ID
		}

		if errors[i].Field != errors[j].Field {
			return errors[i].Field < errors[j].Field
		}

		return errors[i].Error < errors[j].Error
	})
}
