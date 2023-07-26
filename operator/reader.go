package operator

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/modbus"
)

type Reader struct {
	tcpServer *modbus.TCPServer
}

func NewReader(tcpServer *modbus.TCPServer) *Reader {
	reader := &Reader{
		tcpServer: tcpServer,
	}
	return reader
}

func (r *Reader) Trig(address string, targetValue uint32, triggerType data.TriggerType) bool {
	r.tcpServer.Read("DS200", 120)
	//if realtimeValue, has := r.tcpServer.RealtimeValueMap[address]; has {
	if realtimeValue, ok := r.tcpServer.GetRealtimeValue(address); ok {
		switch triggerType {
		case data.LARGER_THAN_TARGET:
			if realtimeValue >= targetValue {
				return true
			} else {
				//logger.Log.Println(address, realtimeValue)
				return false
			}
		case data.EQUAL_TO_TARGET:
			if realtimeValue == targetValue {
				return true
			} else {
				//logger.Log.Println(address, realtimeValue)
				return false
			}
		default:
			logger.Log.Println("wrong trig type")
			return false
		}
	} else {
		return false
	}
}
