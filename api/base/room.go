package base

import (
	"errors"
	"fmt"
	"regexp"
)

// Room - a representation of a room containing a TEC Pi system.
type Room struct {
	ID            string                 `json:"_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Configuration RoomConfiguration      `json:"configuration"`
	Designation   string                 `json:"designation"`
	Devices       []Device               `json:"devices,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

var roomValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,})-[A-z,0-9]+`)

// Validate checks to make sure that the Room's values are valid.
func (r *Room) Validate() error {
	vals := roomValidationRegex.FindStringSubmatch(r.ID)
	if len(vals) == 0 {
		return errors.New("invalid room: _id must match `([A-z,0-9]{2,})-[A-z,0-9]+`")
	}

	if len(r.Name) == 0 {
		return errors.New("invalid room: missing name")
	}

	if len(r.Designation) == 0 {
		return errors.New("invalid room: missing designation")
	}

	if err := r.Configuration.Validate(false); err != nil {
		return fmt.Errorf("invalid room: %s", err)
	}

	return nil
}

// RoomConfiguration - a representation of the configuration of a room.
type RoomConfiguration struct {
	ID          string      `json:"_id"`
	Evaluators  []Evaluator `json:"evaluators,omitempty"`
	Description string      `json:"description,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
}

// Validate checks to make sure that the RoomConfiguration's values are valid.
func (rc *RoomConfiguration) Validate(deepCheck bool) error {
	if len(rc.ID) == 0 {
		return errors.New("invalid room configuration: missing _id")
	}

	if deepCheck {
		if len(rc.Evaluators) == 0 {
			return errors.New("invalid room configuration: at least one evaluator is required")
		}

		for _, evaluator := range rc.Evaluators {
			if err := evaluator.Validate(); err != nil {
				return fmt.Errorf("invalid room configuration: %s", err)
			}
		}
	}

	return nil
}

// Evaluator - a representation of a priority evaluator.
type Evaluator struct {
	ID          string   `json:"_id"`
	CodeKey     string   `json:"codekey,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the Evaluator's values are valid.
func (e *Evaluator) Validate() error {
	if len(e.ID) == 0 {
		return errors.New("invalid evaluator: missing evaluator _id")
	}

	if len(e.CodeKey) == 0 {
		return errors.New("invalid evaluator: missing codekey")
	}

	// default priority to 1000
	if e.Priority == 0 {
		e.Priority = 1000
	}

	return nil
}
