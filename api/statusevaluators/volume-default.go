package statusevaluators

import (
	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/common/log"
)

// VolumeDefaultEvaluator is a constant variable for the name of the evaluator.
const VolumeDefaultEvaluator = "STATUS_VolumeDefault"

// VolumeDefaultCommand is a constant variable for the name of the command.
const VolumeDefaultCommand = "STATUS_Volume"

// VolumeDefault implements the StatusEvaluator struct.
type VolumeDefault struct {
}

// GenerateCommands generates a list of commands for the given devices.
func (p *VolumeDefault) GenerateCommands(room base.Room) ([]StatusCommand, int, error) {
	return generateStandardStatusCommand(room.Devices, VolumeDefaultEvaluator, VolumeDefaultCommand)
}

// EvaluateResponse processes the response information that is given.
func (p *VolumeDefault) EvaluateResponse(room base.Room, label string, value interface{}, Source base.Device, dest base.DestinationDevice) (string, interface{}, error) {
	log.L.Infof("[statusevals] Evaluating response: %s, %s in evaluator %v", label, value, VolumeDefaultCommand)
	return label, value, nil
}
