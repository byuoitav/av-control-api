package couch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	l "github.com/byuoitav/common/log"

	"github.com/byuoitav/common/nerr"
)

type couchReplicationState struct {
	Database    string    `json:"database"`
	DocID       string    `json:"doc_id"`
	ID          string    `json:"id"`
	Source      string    `json:"source"`
	Target      string    `json:"target"`
	State       string    `json:"state"`
	ErrorCount  int       `json:"error_count"`
	StartTime   time.Time `json:"start_time"`
	LastUpdated time.Time `json:"last_updated"`
}

type couchReplicationPayload struct {
	ID           string      `json:"_id"`
	Rev          string      `json:"_rev,omitempty"`
	Source       string      `json:"source"`
	Target       string      `json:"target"`
	CreateTarget bool        `json:"create_target"`
	Continuous   bool        `json:"continous"`
	Selector     interface{} `json:"selector,omitempty"`
}

//Simply returns the replication state.
func (c *CouchDB) GetStatus() (string, error) {
	//check the state of the devices index to see if it's replication or ready.
	state, err := c.CheckReplication("auto_devices")
	if err != nil {
		return "not-ready", err
	}

	return state, nil
}

func (c *CouchDB) CheckReplication(replID string) (string, *nerr.E) {
	l.L.Debugf("Checking to see if replication document %v is already scheduled", replID)

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/_scheduler/docs/_replicator/%v", c.address, replID), nil)
	if err != nil {
		return "", nerr.Translate(err).Addf("Couldn't create request to check replication of %v", replID)
	}

	req.SetBasicAuth(c.username, c.password)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nerr.Translate(err).Addf("Couldn't make request to check replication of %v", replID)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nerr.Translate(err).Addf("Couldn't read error response from couch server while checking for replication %v", replID)
	}

	if resp.StatusCode/100 != 2 {
		ce := CouchError{}
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return "", nerr.Translate(err).Addf("Couldn't Unmarshal response from couch server while checking for replication job %v", replID)
		}

		err = CheckCouchErrors(ce)
		if _, ok := err.(*NotFound); resp.StatusCode == 404 && ok {
			return "not_started", nil
		}

		return "", nerr.Translate(err).Addf("Issue checking replication status of %v", replID)
	}

	//if it's a 200 response, lets see what the state is
	state := couchReplicationState{}
	err = json.Unmarshal(b, &state)
	if err != nil {
		return "", nerr.Translate(err).Addf("Couldn't unmarshal the replication state of %v", replID)
	}
	switch state.State {
	case "completed":
		return "completed", nil
	case "running":
		return "running", nil
	case "started":
		return "started", nil
	case "added":
		return "added", nil
	case "failed":
		return "failed", nil
	case "crashed":
		return "crashed", nil
	default:
		l.L.Errorf("Replication state for %v is in a bad state %v", replID, state.State)
		return state.State, nerr.Create(fmt.Sprintf("Replication of %v is in state %v", replID, state.State), "couch-repl-error")
	}
}
