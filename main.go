package main

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/modbus"
	"cook-robot-controller-go/operator"
	"time"
)

func main() {
	tcpServer := modbus.NewTCPServer("192.168.0.51", 502)
	//tcpServer := modbus.NewTCPServer("127.0.0.1", 502)
	go tcpServer.Run()

	//rpcServer := grpc.NewRPCServer()
	//go rpcServer.Run()
	//
	writer := operator.NewWriter(tcpServer)
	reader := operator.NewReader(tcpServer)
	controller := core.NewController(writer, reader)
	go controller.Run()
	//
	delayAction1 := action.NewDelayAction(1000)
	//delayAction2 := action.NewDelayAction(1500)
	//controller.AddAction(delayAction1)
	//controller.AddAction(delayAction2)
	//
	//controlAction1 := action.NewTemperatureControlAction(
	//	data.TEMPERATURE_CONTROL_WORD_ADDRESS,
	//	data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_TARGET_ADDRESS, 1001))
	//controller.AddAction(controlAction1)
	//triggerAction := action.NewTriggerAction(data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, 0))
	//controller.AddAction(triggerAction)
	////controller.AddAction(delayAction2)
	////
	time.Sleep(time.Second)

	controlAction1 := action.NewPumpControlAction(data.PUMP_CONTROL_WORD_ADDRESS, data.PUMP_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.PUMP_NUMBER_ADDRESS, 1), data.NewAddressValue(data.PUMP_DURATION_ADDRESS, 8000))

	//controlAction2 := action.NewAxleLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 100), data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, 1000))
	controlAction3 := action.NewPumpControlAction(data.PUMP_CONTROL_WORD_ADDRESS, data.PUMP_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.PUMP_NUMBER_ADDRESS, 2), data.NewAddressValue(data.PUMP_DURATION_ADDRESS, 4000))
	controller.AddAction(controlAction1)
	controller.AddAction(controlAction3)
	//controller.AddAction(controlAction2)
	controller.AddAction(delayAction1)

	go controller.Start()

	for true {
		time.Sleep(1000 * time.Millisecond)
	}
}
