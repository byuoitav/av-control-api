package couch

import (
	"errors"
	"fmt"

	"github.com/byuoitav/av-control-api/api/base"
)

type room struct {
	Rev string `json:"_rev,omitempty"`
	*base.Room
}

type roomConfiguration struct {
	Rev string `json:"_rev,omitempty"`
	*base.RoomConfiguration
}

func (c *CouchDB) GetRoom(id string) (base.Room, error) {
	room, err := c.getRoom(id)
	if err != nil {
		return base.Room{}, err
	}
	//if err was nil then room may be.
	return *room.Room, nil
}

func (c *CouchDB) getRoom(id string) (room, error) {
	var toReturn room

	// get the base room
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%v", ROOMS, id), "", nil, &toReturn)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get room %s: %s", id, err))
	}

	// get the devices in room
	devices, err := c.GetDevicesByRoom(id)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get devices in room %s: %s", id, err))
	}

	// fill devices into room
	for _, device := range devices {
		toReturn.Devices = append(toReturn.Devices, device)
	}

	// get room configuration
	toReturn.Configuration, err = c.GetRoomConfiguration(toReturn.Configuration.ID)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to get room configuration %s for room %s: %s", toReturn.Configuration.ID, id, err))
	}

	return toReturn, nil
}

func (c *CouchDB) GetRoomConfiguration(id string) (base.RoomConfiguration, error) {
	rc, err := c.getRoomConfiguration(id)
	switch {
	case err != nil:
		return base.RoomConfiguration{}, err
	case rc.RoomConfiguration == nil:
		return base.RoomConfiguration{}, fmt.Errorf("no room configuration %q found", id)
	}

	return *rc.RoomConfiguration, err
}

func (c *CouchDB) getRoomConfiguration(id string) (roomConfiguration, error) {
	var toReturn roomConfiguration

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", ROOM_CONFIGURATIONS, id), "", nil, &toReturn)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get room configuration %s: %s", id, err))
	}

	return toReturn, err
}
