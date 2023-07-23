package main

import (
	"cook-robot-controller-go/config"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/grpc"
	"cook-robot-controller-go/instruction"
	"cook-robot-controller-go/modbus"
	"cook-robot-controller-go/operator"
	"time"
)

func main() {
	tcpServer := modbus.NewTCPServer(config.App.Modbus.TargetHost, config.App.Modbus.TargetPort, config.App.DebugMode)
	go tcpServer.Run()

	//
	writer := operator.NewWriter(tcpServer)
	reader := operator.NewReader(tcpServer)
	controller := core.NewController(writer, reader, tcpServer, config.App.DebugMode)
	go controller.Run()

	rpcServer := grpc.NewGRPCServer(config.App.GRPC.Host, config.App.GRPC.Port, controller)
	go rpcServer.Run()

	time.Sleep(500 * time.Millisecond)

	resetXYTInstruction := instruction.NewResetXYTInstruction()
	controller.CurrentCommandName = "reset"
	resetXYTInstruction.AddToController(controller)
	go controller.Start()
	//
	//delayAction1 := action.NewDelayAction(1000)
	//delayAction2 := action.NewDelayAction(1500)
	//controller.AddAction(delayAction1)
	//controller.AddAction(delayAction2)
	//
	//controlAction1 := action.NewTemperatureControlAction(
	//	data.TEMPERATURE_CONTROL_WORD_ADDRESS,
	//	data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_ADDRESS, 1001))
	//controller.AddAction(controlAction1)
	//triggerAction := action.NewTriggerAction(data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, 0))
	//controller.AddAction(triggerAction)
	////controller.AddAction(delayAction2)
	////
	//instruction.NewResetXYInstruction().AddToController(controller)
	//instruction.NewSeasoningInstruction(map[string]uint32{"5": 10, "9": 10}).AddToController(controller)
	//instruction.NewResetRTInstruction().AddToController(controller)
	//controller.AddAction(action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS, data.Y_RESET_STATUS_WORD_ADDRESS))
	//for {
	//	//controller.AddAction(action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS, data.X_RESET_STATUS_WORD_ADDRESS))
	//	controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3500), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//	controller.AddAction(action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
	//		data.R1_ROTATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
	//		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 600),
	//		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 2)))
	//	controller.AddAction(action.NewDelayAction(5000))
	//	controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3300), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//	controller.AddAction(action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
	//		data.R1_ROTATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 2),
	//		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 800),
	//		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 2)))
	//	controller.AddAction(action.NewDelayAction(5000))
	//	controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3000), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//	controller.AddAction(action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
	//		data.R1_ROTATE_STATUS_WORD_ADDRESS,
	//		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 3),
	//		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 1200),
	//		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 3)))
	//	controller.AddAction(action.NewDelayAction(10000))
	//
	//	controller.Start()
	//	time.Sleep(5 * time.Second)
	//}

	//
	//pumpAction1 := action.NewPumpControlAction(data.PUMP_1_CONTROL_WORD_ADDRESS, data.PUMP_1_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_1_DURATION_ADDRESS, 2000))
	//pumpAction2 := action.NewPumpControlAction(data.PUMP_2_CONTROL_WORD_ADDRESS, data.PUMP_2_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_2_DURATION_ADDRESS, 4000))
	//pumpAction3 := action.NewPumpControlAction(data.PUMP_6_CONTROL_WORD_ADDRESS, data.PUMP_6_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_6_DURATION_ADDRESS, 6000))
	//groupAction := action.NewGroupAction()
	//groupAction.AddAction(pumpAction1)
	//groupAction.AddAction(pumpAction2)
	//groupAction.AddAction(pumpAction3)
	//
	//controller.AddAction(groupAction)
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3500), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))

	//controller.AddAction(action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS, data.X_RESET_STATUS_WORD_ADDRESS))
	//controller.AddAction(action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS, data.Y_RESET_STATUS_WORD_ADDRESS))
	//xPositionList := []uint32{0, 54396, 42690, 25937, 9069, 20}
	//yPositionList := []uint32{0, 2000, 3200, 3485, 3500, 1759, 10, 4200, 4700}
	//time.Sleep(time.Second)
	//rand.Seed(time.Now().UnixNano())
	//for i := 0; i < 10; i++ {
	//	// 生成一个随机索引
	//	xIndex := rand.Intn(len(xPositionList))
	//yIndex := rand.Intn(len(yPositionList))

	// 根据随机索引从数组中获取数值
	//xPosition := xPositionList[xIndex]
	//yPosition := yPositionList[yIndex]
	//controller.AddAction(action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, xPosition), data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 20000)))
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, yPosition), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//}

	//controller.AddAction(action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, 54396), data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 20000)))
	//
	//controller.AddAction(action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS, data.X_RESET_STATUS_WORD_ADDRESS))
	//controller.AddAction(action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS, data.Y_RESET_STATUS_WORD_ADDRESS))
	//
	//controller.AddAction(action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, 54396), data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 20000)))
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3485), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//controller.AddAction(action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS, data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, 2)))
	//controller.AddAction(action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS, data.R1_ROTATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 3), data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 500),
	//	data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 5)))
	//
	//controller.AddAction(action.NewDelayAction(5000))
	//
	//controller.AddAction(action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS, data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_ADDRESS, 150)))
	//
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3200), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//
	//controller.AddAction(action.NewPumpControlAction(data.PUMP_CONTROL_WORD_ADDRESS, data.PUMP_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_NUMBER_ADDRESS, 6), data.NewAddressValue(data.PUMP_DURATION_ADDRESS, 5000)))
	//
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 2000), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//
	//controller.AddAction(action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS, data.X_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, 42690), data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, 20000)))
	//
	//controller.AddAction(action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS, data.SHAKE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
	//	data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000), data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
	//	data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 2000)))
	//
	//controller.AddAction(action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, 3485), data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, 600)))
	//
	//controller.AddAction(action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS, data.R1_ROTATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 3), data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 1500),
	//	data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 5)))
	//
	//controller.AddAction(action.NewDelayAction(30000))
	//
	//controller.AddAction(action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS, data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0)))
	//controller.AddAction(action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS, data.R1_ROTATE_STATUS_WORD_ADDRESS))
	//controller.AddAction(action.NewStopAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS, data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS))
	//
	//controller.AddAction(action.NewDelayAction(2000))
	//
	//controller.AddAction(action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS, data.DOOR_UNLOCK_STATUS_WORD_ADDRESS))
	//
	//controller.AddAction(action.NewDishOutControlAction(data.DISH_OUT_CONTROL_WORD_ADDRESS, data.DISH_OUT_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.DISH_OUT_AMOUNT_ADDRESS, 3),
	//	data.NewAddressValue(data.DISH_OUT_UPWARD_SPEED_ADDRESS, 400), data.NewAddressValue(data.DISH_OUT_DOWNWARD_SPEED_ADDRESS, 800),
	//	data.NewAddressValue(data.DISH_OUT_UPWARD_POSITION_ADDRESS, 4200), data.NewAddressValue(data.DISH_OUT_DOWNWARD_POSITION_ADDRESS, 4700)))
	//controlAction3 := action.NewPumpControlAction(data.PUMP_CONTROL_WORD_ADDRESS, data.PUMP_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_NUMBER_ADDRESS, 2), data.NewAddressValue(data.PUMP_DURATION_ADDRESS, 4000))
	//controlAction4 := action.NewDishOutControlAction(data.DISH_OUT_CONTROL_WORD_ADDRESS, data.DISH_OUT_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.DISH_OUT_AMOUNT_ADDRESS, 3),
	//	data.NewAddressValue(data.DISH_OUT_UPWARD_SPEED_ADDRESS, 400),
	//	data.NewAddressValue(data.DISH_OUT_DOWNWARD_SPEED_ADDRESS, 800),
	//	data.NewAddressValue(data.DISH_OUT_UPWARD_POSITION_ADDRESS, 2000),
	//	data.NewAddressValue(data.DISH_OUT_DOWNWARD_POSITION_ADDRESS, 4000))

	//controlAction5 := action.NewPumpControlAction(data.PUMP_CONTROL_WORD_ADDRESS, data.PUMP_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.PUMP_NUMBER_ADDRESS, 3),
	//	data.NewAddressValue(data.PUMP_DURATION_ADDRESS, 4000))

	//controller.AddAction(controlAction1)
	//controller.AddAction(controlAction3)
	//controller.AddAction(controlAction2)
	//controller.AddAction(delayAction1)
	//controller.AddAction(controlAction4)
	//controller.AddAction(controlAction5)

	//go controller.Start()

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

// $env:CGO_ENABLED="0"
//  $env:GOOS="linux"
//$env:GOARCH="arm"
