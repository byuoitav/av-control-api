package couch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
)

/* Consts */
const (
	BUILDINGS           = "buildings"
	ROOMS               = "rooms"
	ROOM_ATTACHMENTS    = "room_attachments"
	DEVICES             = "devices"
	DEVICE_STATES       = "device-state"
	DEVICE_TYPES        = "device-types"
	ROOM_CONFIGURATIONS = "room_configurations"
	UI_CONFIGS          = "ui-configuration"
	OPTIONS             = "options"
	ICONS               = "Icons"
	ROLES               = "DeviceRoles"
	ROOM_DESIGNATIONS   = "RoomDesignations"
	CLOSURE_CODES       = "ClosureCodes"
	TAGS                = "Tags"
	DMPSLIST            = "dmps"
	LAB_CONFIGS         = "lab-attendance-config"
	SCHEDULING_CONFIGS  = "scheduling-configs"

	deviceMonitoring = "device-monitoring"
	MENUTREE         = "MenuTree"
	ATTRIBUTES       = "attributes"

	DEPLOY = "deployment-information"
	CAMPUS = "campus-deployment-info"
)

// CouchDB .
type CouchDB struct {
	address  string
	username string
	password string

	IgnoreReadyChecks bool
}

// NewDB .
func NewDB(address, username, password string) *CouchDB {
	return &CouchDB{
		address:  strings.Trim(address, "/"),
		username: username,
		password: password,
	}
}

func (c *CouchDB) req(method, endpoint, contentType string, body []byte) (string, []byte, error) {
	errMsg := "unable to make request against couch"

	if len(c.address) == 0 {
		return "", nil, fmt.Errorf("%s: couch address not set", errMsg)
	}

	url := fmt.Sprintf("%s/%s", c.address, endpoint)
	url = strings.TrimSpace(url)

	// build request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return "", nil, fmt.Errorf("%s: %s", errMsg, err)
	}

	// add auth
	if len(c.username) > 0 && len(c.password) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	// add headers
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// validate that couch is ready, wait if it isn't
	if !c.IgnoreReadyChecks {
		c.waitUntilReady()
	}

	// execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %s", errMsg, err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %s", errMsg, err)
	}

	if resp.StatusCode/100 != 2 {
		ce := CouchError{}
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return "", nil, fmt.Errorf("%s: received a non-200 response from %s. body: %s", errMsg, url, b)
		}

		return "", nil, CheckCouchErrors(ce)
	}

	return resp.Header.Get("content-type"), b, nil
}

var (
	checkReady sync.Once
	readyMu    sync.Mutex
)

func (c *CouchDB) waitUntilReady() {
	checkReady.Do(func() {
		readyMu.Lock()
		defer readyMu.Unlock()

		// +deployment not-required
		for len(os.Getenv("STOP_REPLICATION")) == 0 {
			// wait until database is ready
			state, err := c.GetStatus()
			if err != nil || state != "completed" {
				log.L.Warnf("Database replication in state %v (error: %s); Retrying in 5 seconds", state, err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.L.Infof("Database replication in state %v. Allowing CouchDB requests now.", state)
			break
		}
	})

	readyMu.Lock()
	readyMu.Unlock()
}

// MakeRequest .
func (c *CouchDB) MakeRequest(method, endpoint, contentType string, body []byte, toFill interface{}) error {
	respType, body, err := c.req(method, endpoint, contentType, body)
	if err != nil {
		return err
	}

	if !strings.EqualFold(respType, "application/json") {
		return fmt.Errorf("unexpected response content-type: expected %s, but got %s", "application/json", respType)
	}

	if toFill == nil {
		return nil
	}

	// unmarshal body
	err = json.Unmarshal(body, toFill)
	if err != nil {
		return fmt.Errorf("unable to make request against couch: %s", err)
	}

	return nil
}

func (c *CouchDB) ExecuteQuery(query IDPrefixQuery, responseToFill interface{}) error {
	//	var toFill interface{}

	// marshal query
	b, err := json.Marshal(query)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to marshal query: %s", err))
	}

	//var toReturn []interface{}
	var database string
	//	var sliceType reflect.Type

	switch responseToFill.(type) {
	case buildingQueryResponse:
		database = BUILDINGS
		//	sliceType = reflect.TypeOf(responseToFill)
	}

	// execute query
	err = c.MakeRequest("POST", fmt.Sprintf("%s/find", database), "application/json", b, &responseToFill)
	if err != nil {
		return err
	}

	//	sliceType = reflect.ValueOf(responseToFill)

	return nil
}

func CheckCouchErrors(ce CouchError) error {
	switch strings.ToLower(ce.Error) {
	case "not_found":
		return &NotFound{fmt.Sprintf("The ID requested was unknown. Message: %v.", ce.Reason)}
	case "conflict":
		return &Conflict{fmt.Sprintf("There was a conflict updating/creating the document: %v", ce.Reason)}
	case "bad_request":
		return &BadRequest{fmt.Sprintf("The request was bad: %v", ce.Reason)}
	default:
		return errors.New(fmt.Sprintf("unknown error type: %v. Message: %v", ce.Error, ce.Reason))
	}
}

type IDPrefixQuery struct {
	Selector struct {
		ID struct {
			GT    string `json:"$gt,omitempty"`
			LT    string `json:"$lt,omitempty"`
			Regex string `json:"$regex,omitempty"`
		} `json:"_id"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

type CouchUpsertResponse struct {
	OK  bool   `json:"ok"`
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

type CouchError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type NotFound struct {
	msg string
}

func (n NotFound) Error() string {
	return n.msg
}

type Conflict struct {
	msg string
}

func (c Conflict) Error() string {
	return c.msg
}

type BadRequest struct {
	msg string
}

func (br BadRequest) Error() string {
	return br.msg
}
