package statusevaluators

import (
	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/common/log"
)

// BlankedDefaultEvaluator is a constant variable for the name of the evaluator.
const BlankedDefaultEvaluator = "STATUS_BlankedDefault"

// BlankedDefaultCommand is a constant variable for the name of the command.
const BlankedDefaultCommand = "STATUS_Blanked"

// BlankedDefault implements the StatusEvaluator struct.
type BlankedDefault struct {
}

// GenerateCommands generates a list of commands for the given devices.
func (p *BlankedDefault) GenerateCommands(room base.Room) ([]StatusCommand, int, error) {
	return generateStandardStatusCommand(room.Devices, BlankedDefaultEvaluator, BlankedDefaultCommand)
}

// EvaluateResponse is supposed to evaluate the response...but it seems like it's just logging a statement...
func (p *BlankedDefault) EvaluateResponse(room base.Room, label string, value interface{}, Source base.Device, dest base.DestinationDevice) (string, interface{}, error) {
	log.L.Infof("[statusevals] Evaluating response: %s, %s in evaluator %v", label, value, BlankedDefaultEvaluator)
	return label, value, nil
}
