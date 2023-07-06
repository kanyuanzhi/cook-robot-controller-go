package operator

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/modbus"
)

type Reader struct {
	tcpServer        *modbus.TCPServer
	realtimeValueMap map[string]uint32
}

func NewReader(tcpServer *modbus.TCPServer) *Reader {
	realtimeValueMap := make(map[string]uint32)

	realtimeValueMap[data.X_RESET_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.Y_RESET_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.Z_RESET_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.R1_RESET_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.R2_RESET_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.X_LOCATE_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.Y_LOCATE_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.Z_LOCATE_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.R1_LOCATE_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.R2_LOCATE_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.R1_ROTATE_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.DISH_OUT_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.PUMP_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.DOOR_LOCK_STATUS_WORD_ADDRESS] = 100

	realtimeValueMap[data.TEMPERATURE_STATUS_WORD_ADDRESS] = 100
	realtimeValueMap[data.TEMPERATURE_BOTTOM_ADDRESS] = 0
	realtimeValueMap[data.TEMPERATURE_INFRARED_ADDRESS] = 0
	realtimeValueMap[data.TEMPERATURE_WARNING_ADDRESS] = 0

	reader := &Reader{
		tcpServer:        tcpServer,
		realtimeValueMap: realtimeValueMap,
	}
	return reader
}

func (r *Reader) Trig(address string, targetValue uint32) bool {
	r.tcpServer.Read("DS0", 120)
	if realtimeValue, has := r.tcpServer.RealtimeValueMap[address]; has {
		if realtimeValue >= targetValue {
			return true
		} else {
			return false
		}
	} else {
		logger.Log.Println("no address")
		return true
	}
}
