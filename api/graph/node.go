package graph

import (
	"crypto/sha1"
	"math/big"

	"github.com/byuoitav/av-control-api/api"
)

type Node struct {
	id int64
	*api.Device
}

func (n Node) ID() int64 {
	return n.id
}

func (n Node) DOTID() string {
	return string(n.Device.ID)
}

func NodeID(id api.DeviceID) int64 {
	sum := sha1.Sum([]byte(id))

	var i big.Int
	i.SetBytes(sum[:])
	return i.Int64()
}
