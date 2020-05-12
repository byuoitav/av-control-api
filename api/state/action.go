package state

import (
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type action struct {
	ID    api.DeviceID
	Req   *http.Request
	Order *int

	//Response <-struct{
	//	respBody
	//	Headers
	//	error: Do
	//}
}
