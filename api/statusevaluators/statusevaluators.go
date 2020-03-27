package statusevaluators

import (
	"strings"

	"github.com/byuoitav/av-control-api/api/db"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/common/log"
)

// StatusEvaluator defines the common functions for all StatusEvaluators.
type StatusEvaluator interface {
	//Generates action list
	GenerateCommands(room base.Room) ([]StatusCommand, int, error)

	//Evaluate Response
	EvaluateResponse(room base.Room, label string, value interface{}, Source base.Device, Destination base.DestinationDevice) (string, interface{}, error)
}

// StatusEvaluatorMap is a map of the different StatusEvaluators used.
//TODO: we shoud grab the keys from constants in the evaluators themselves
var StatusEvaluatorMap = map[string]StatusEvaluator{
	"STATUS_PowerDefault":       &PowerDefault{},
	"STATUS_BlankedDefault":     &BlankedDefault{},
	"STATUS_MutedDefault":       &MutedDefault{},
	"STATUS_InputDefault":       &InputDefault{},
	"STATUS_VolumeDefault":      &VolumeDefault{},
	"STATUS_InputVideoSwitcher": &InputVideoSwitcher{},
	"STATUS_InputDSP":           &InputDSP{},
	"STATUS_MutedDSP":           &MutedDSP{},
	"STATUS_VolumeDSP":          &VolumeDSP{},
	"STATUS_Tiered_Switching":   &InputTieredSwitcher{},
}

func generateStandardStatusCommand(devices []base.Device, evaluatorName string, commandName string) ([]StatusCommand, int, error) {
	var count int

	log.L.Infof("[statusevals] Generating status commands from %v", evaluatorName)
	var output []StatusCommand

	//iterate over each device
	for _, device := range devices {

		log.L.Infof("[statusevals] Considering device: %s", device.Name)

		for id, command := range device.Type.Commands {
			if strings.HasPrefix(id, FLAG) && strings.EqualFold(id, commandName) {
				log.L.Info("[statusevals] Command found")

				//every power command needs an address parameter
				parameters := make(map[string]string)
				parameters["address"] = device.Address

				//build destination device
				var destinationDevice base.DestinationDevice
				for _, role := range device.Roles {
					if role.ID == "AudioOut" {
						destinationDevice.AudioDevice = true
					}

					if role.ID == "VideoOut" {
						destinationDevice.Display = true
					}
				}

				destinationDevice.Device = device

				log.L.Infof("[statusevals] Adding command: %s to action list with device %s", id, device.ID)
				output = append(output, StatusCommand{
					ActionID:          id,
					Action:            command,
					Device:            device,
					Parameters:        parameters,
					DestinationDevice: destinationDevice,
					Generator:         evaluatorName,
				})
				count++

				////////////////////////
				///// MIRROR STUFF /////
				if base.HasRole(device, "MirrorMaster") {
					for _, port := range device.Ports {
						if port.ID == "mirror" {
							DX, err := db.GetDB().GetDevice(port.DestinationDevice)
							if err != nil {
								return output, count, err
							}

							_, err = DX.GetCommandByID(commandName)
							if err != nil {
								continue
							}

							destinationDevice.Device = DX

							log.L.Infof("[statusevals] Adding command: %s to action list with device %s", id, DX.ID)
							output = append(output, StatusCommand{
								ActionID:          id,
								Action:            command,
								Device:            DX,
								Parameters:        parameters,
								DestinationDevice: destinationDevice,
								Generator:         evaluatorName,
							})
							count++
						}
					}
				}
				///// MIRROR STUFF /////
				////////////////////////
			}

		}

	}

	return output, count, nil

}

// FindDevice searches a list of devices for the device specified by the given ID and returns it
func FindDevice(deviceList []base.Device, searchID string) base.Device {
	for i := range deviceList {
		if deviceList[i].ID == searchID {
			return deviceList[i]
		}
	}

	return base.Device{}
}

// FilterDevicesByRole searches a list of devices for the devices that have the given roles, and returns a new list of those devices
func FilterDevicesByRole(deviceList []base.Device, roleID string) []base.Device {
	var toReturn []base.Device

	for _, device := range deviceList {
		if device.HasRole(roleID) {
			toReturn = append(toReturn, device)
		}
	}

	return toReturn
}
