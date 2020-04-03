package state

import (
	"fmt"

	"github.com/byuoitav/av-control-api/api/db"
	"github.com/byuoitav/av-control-api/api/rest"
	"github.com/byuoitav/av-control-api/api/statusevaluators"
	"github.com/byuoitav/common/log"
	"github.com/fatih/color"
)

//GetRoomState assesses the state of the room and returns a PublicRoom object.
func GetRoomState(building, roomName, env string) (rest.PublicRoom, error) {

	color.Set(color.FgHiCyan, color.Bold)
	log.L.Info("[state] getting room state...")
	color.Unset()

	roomID := fmt.Sprintf("%v-%v", building, roomName)
	room, err := db.GetDB().GetRoom(roomID)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	//we get the number of actions generated
	commands, count, err := GenerateStatusCommands(room, statusevaluators.StatusEvaluatorMap)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	responses, err := RunStatusCommands(commands, env)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	roomStatus, err := EvaluateResponses(room, responses, count)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	roomStatus.Building = building
	roomStatus.Room = roomName

	color.Set(color.FgHiGreen, color.Bold)
	log.L.Info("[state] successfully retrieved room state")
	color.Unset()

	return roomStatus, nil
}

//SetRoomState changes the state of the room and returns a PublicRoom object.
func SetRoomState(target rest.PublicRoom, env, requestor string) (rest.PublicRoom, error) {
	log.L.Infof("Requestor: %v\n", requestor)
	log.L.Infof("%s", color.HiBlueString("[state] setting room state..."))

	roomID := fmt.Sprintf("%v-%v", target.Building, target.Room)
	room, err := db.GetDB().GetRoom(roomID)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	//so here we need to know how many things we're actually expecting.
	actions, count, err := GenerateActions(room, target, requestor)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	responses, err := ExecuteActions(actions, env, requestor)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	//here's where we then pass that information through so that we can make a decent decision.
	report, err := EvaluateResponses(room, responses, count)
	if err != nil {
		return rest.PublicRoom{}, err
	}

	report.Building = target.Building
	report.Room = target.Room

	color.Set(color.FgHiGreen, color.Bold)
	log.L.Info("[state] successfully set room state")
	color.Unset()

	return report, nil
}
