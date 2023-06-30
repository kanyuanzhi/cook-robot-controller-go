package main

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/operator"
	"time"
)

func main() {
	writer := operator.NewWriter()
	reader := operator.NewReader()
	controller := core.NewController(writer, reader)
	go writer.Run()
	go controller.Run()

	delayAction1 := action.NewDelayAction(1000)
	delayAction2 := action.NewDelayAction(1500)
	controller.AddAction(delayAction1)
	controller.AddAction(delayAction2)

	controlAction1 := action.NewTemperatureControlAction(
		data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_TARGET_ADDRESS, 1001))
	controller.AddAction(controlAction1)
	triggerAction := action.NewTriggerAction(data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, 0))
	controller.AddAction(triggerAction)
	//controller.AddAction(delayAction2)
	//
	controlAction2 := action.NewAxleLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 100), data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, 1000))
	//controller.AddAction(controlAction1)
	controller.AddAction(controlAction2)
	go controller.Start()

	for true {
		time.Sleep(1000 * time.Millisecond)
	}
}
