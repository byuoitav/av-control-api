package db

import (
	"os"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/av-control-api/api/db/couch"
	"github.com/byuoitav/common/log"
)

// DB .
type DB interface {
	/* crud functions */
	// Room
	GetRoom(id string) (base.Room, error)

	// device
	GetDevice(id string) (base.Device, error)

	/* bulk functions */
	GetAllBuildings() ([]base.Building, error)
}

var address string
var username string
var password string

var database DB

func init() {
	address = os.Getenv("DB_ADDRESS")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
}

// GetDB returns the instance of the database to use.
func GetDB() DB {
	if len(address) == 0 {
		log.L.Errorf("DB_ADDRESS is not set.")
	}

	// TODO add logic to "pick" which db to create
	if database == nil {
		database = couch.NewDB(address, username, password)
	}

	return database
}

// GetDBWithCustomAuth returns an instance of the database with a custom authentication
func GetDBWithCustomAuth(address, username, password string) DB {
	return couch.NewDB(address, username, password)
}
