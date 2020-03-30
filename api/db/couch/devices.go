package couch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/av-control-api/api/base"
)

type device struct {
	Rev string `json:"_rev,omitempty"`
	*base.Device
}

type deviceType struct {
	Rev string `json:"_rev,omitempty"`
	*base.DeviceType
}

type deviceQueryResponse struct {
	Docs     []device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type deviceTypeQueryResponse struct {
	Docs     []deviceType `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}

// GetDevice .
func (c *CouchDB) GetDevice(id string) (base.Device, error) {
	device, err := c.getDevice(id)
	switch {
	case err != nil:
		return base.Device{}, err
	case device.Device == nil:
		return base.Device{}, fmt.Errorf("device not found")
	default:
		return *device.Device, err
	}
}

func (c *CouchDB) getDevice(id string) (device, error) {
	var toReturn device // get the device
	err := c.MakeRequest("GET", fmt.Sprintf("%s/%v", DEVICES, id), "", nil, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get device %s: %s", id, err)
	}

	if len(toReturn.ID) == 0 {
		return toReturn, fmt.Errorf("failed to get device %s: %s", id, err)
	}

	// get its device type
	toReturn.Type, err = c.getDeviceType(toReturn.Type.ID)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get device type (%s) to get device %s: %s", toReturn.Type.ID, id, err)
	}

	return toReturn, err
}

func (c *CouchDB) getDeviceType(id string) (base.DeviceType, error) {
	var toReturn deviceType

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", DEVICE_TYPES, id), "", nil, &toReturn)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get device type %s: %s", id, err))
	}

	return *toReturn.DeviceType, err
}

// GetDevicesByRoom .
func (c *CouchDB) GetDevicesByRoom(roomID string) ([]base.Device, error) {
	var toReturn []base.Device

	devices, err := c.getDevicesByRoom(roomID)
	if err != nil {
		return toReturn, err
	}

	for _, device := range devices {
		toReturn = append(toReturn, *device.Device)
	}

	return toReturn, nil
}

func (c *CouchDB) getDevicesByRoom(roomID string) ([]device, error) {
	var toReturn []device

	// create query
	var query IDPrefixQuery
	query.Selector.ID.GT = fmt.Sprintf("%v-", roomID)
	query.Selector.ID.LT = fmt.Sprintf("%v.", roomID)
	query.Limit = 1000

	// query devices
	toReturn, err := c.getDevicesByQuery(query, true)
	if err != nil {
		return toReturn, fmt.Errorf("failed getting devices in room %s: %s", roomID, err)
	}

	return toReturn, nil
}

func (c *CouchDB) getDevicesByQuery(query IDPrefixQuery, includeType bool) ([]device, error) {
	var toReturn []device

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal devices query: %s", err)
	}

	// make query for devices
	var resp deviceQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%s/_find", DEVICES), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to query devices: %s", err)
	}

	if includeType {
		// get all types
		types, err := c.GetAllDeviceTypes()
		if err != nil {
			return toReturn, fmt.Errorf("failed to get devices types for devices query:%s", err)
		}

		// make a map of type.ID -> type
		typesMap := make(map[string]base.DeviceType)
		for _, t := range types {
			typesMap[t.ID] = t
		}

		// fill in device types
		for _, d := range resp.Docs {
			d.Type = typesMap[d.Type.ID]
		}
	}

	// return each document
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}

func (c *CouchDB) getDeviceTypesByQuery(query IDPrefixQuery) ([]deviceType, error) {
	var toReturn []deviceType

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to marshal device types query: %s", err))
	}

	// make query for types
	var resp deviceTypeQueryResponse
	err = c.MakeRequest("POST", fmt.Sprintf("%s/_find", DEVICE_TYPES), "application/json", b, &resp)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed to query device types: %s", err))
	}

	// return each document
	for _, doc := range resp.Docs {
		toReturn = append(toReturn, doc)
	}

	return toReturn, nil
}

func (c *CouchDB) GetAllDeviceTypes() ([]base.DeviceType, error) {
	var toReturn []base.DeviceType

	// create all device types query
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 5000

	// execute query
	types, err := c.getDeviceTypesByQuery(query)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("failed getting all device types: %s", err))
	}

	// return the struct part
	for _, t := range types {
		toReturn = append(toReturn, *t.DeviceType)
	}

	return toReturn, nil
}
