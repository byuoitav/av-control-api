package base

import "errors"

// Building - the representation about a building containing a TEC Pi system.
type Building struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate determines if the current values for the building's attributes are valid or not.
func (b *Building) Validate() error {
	if len(b.ID) < 2 {
		return errors.New("invalid building: id must be at least 2 characters long")
	}
	return nil
}
