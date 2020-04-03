package commandevaluators

import (
	"strings"

	"github.com/byuoitav/av-control-api/api/base"
)

//This file contains common 'helper' functions.

//Checks an action list to see if it has a device (by name, room, and building) already in it,
//if so, it returns the index of the device, if not -1.
func checkActionListForDevice(a []base.ActionStructure, d string, room string, building string) (index int) {
	for i, curDevice := range a {
		if checkDevicesEqual(curDevice.Device, d, room, building) {
			return i
		}
	}
	return -1
}

func checkDevicesEqual(dev base.Device, name string, room string, building string) bool {
	splits := strings.Split(dev.GetDeviceRoomID(), "-")
	return strings.EqualFold(dev.ID, name) &&
		strings.EqualFold(splits[1], room) &&
		strings.EqualFold(splits[0], building)
}

// CheckCommands searches a list of Commands to see if it contains any command by the name given.
// returns T/F, as well as the command if true.
func CheckCommands(commands map[string]base.Command, commandName string) (bool, base.Command) {
	for id, c := range commands {
		if strings.EqualFold(id, commandName) {
			return true, c
		}
	}
	return false, base.Command{}
}

func markAsOverridden(action base.ActionStructure, structs ...[]*base.ActionStructure) {
	for i := 0; i < len(structs); i++ {
		for j := 0; j < len(structs[i]); j++ {
			if structs[i][j].Equals(action) {
				structs[i][j].Overridden = true
			}
		}
	}
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
