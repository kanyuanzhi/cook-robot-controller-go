package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"math"
	"time"
)

func (d DoorUnlockInstruction) ExecuteImmediately(controller *core.Controller) {
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.ExecuteImmediately(doorUnlockControlAction)
}

func (p PauseToAddInstruction) ExecuteImmediately(controller *core.Controller) {
	// y轴处在翻炒位延时时才允许暂停中途加料
	if !controller.IsPausePermitted {
		logger.Log.Println("y轴不在翻炒位，禁止中途加料")
		return
	}
	//controller.IsPausing = true
	//controller.IsPausingWithMovingFinished = false
	controller.Pause()
	//yLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
	//	data.Y_LOCATE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
	//	data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	//controller.ExecuteImmediately(yLocateControlAction)
	//triggerAction := action.NewTriggerAction(data.NewAddressValue(data.Y_LOCATE_STATUS_WORD_ADDRESS, 100), data.EQUAL_TO_TARGET, data.Y_LOCATE_CONTROL_WORD_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION])
	//controller.ExecuteImmediately(triggerAction)
	//controller.IsPausingWithMovingFinished = true
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.ExecuteImmediately(doorUnlockControlAction)
}

func (r ResumeInstruction) ExecuteImmediately(controller *core.Controller) {
	controller.Resume()
}

func (h HeatInstruction) ExecuteImmediately(controller *core.Controller) {
	if h.JudgeType == NO_JUDGE {
		temperatureControlAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
			data.TEMPERATURE_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.TEMPERATURE_ADDRESS, uint32(math.Round(h.Temperature*10))))
		controller.ExecuteImmediately(temperatureControlAction)
	} else {
		logger.Log.Println("只有无控制的加热指令可以立即执行")
	}
}

func (r ShutdownInstruction) ExecuteImmediately(controller *core.Controller) {
	logger.Log.Println("[终止]停止所有动作")
	controller.Shutdown()
	stopActions := make(chan action.Actioner, 100)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)
	stopActions <- axisR1StopAction

	axisXStopAction := action.NewStopAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS)
	stopActions <- axisXStopAction

	axisYStopAction := action.NewStopAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS)
	stopActions <- axisYStopAction

	axisZStopAction := action.NewStopAction(data.Z_LOCATE_CONTROL_WORD_ADDRESS,
		data.Z_LOCATE_STATUS_WORD_ADDRESS)
	stopActions <- axisZStopAction

	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	stopActions <- temperatureResetAction

	dishOutStopAction := action.NewStopAction(data.DISH_OUT_CONTROL_WORD_ADDRESS,
		data.DISH_OUT_STATUS_WORD_ADDRESS)
	stopActions <- dishOutStopAction

	shakeStopAction := action.NewStopAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS)
	stopActions <- shakeStopAction

	pump1StopAction := action.NewStopAction(data.PUMP_1_CONTROL_WORD_ADDRESS,
		data.PUMP_1_STATUS_WORD_ADDRESS)
	stopActions <- pump1StopAction

	pump2StopAction := action.NewStopAction(data.PUMP_2_CONTROL_WORD_ADDRESS,
		data.PUMP_2_STATUS_WORD_ADDRESS)
	stopActions <- pump2StopAction

	pump3StopAction := action.NewStopAction(data.PUMP_3_CONTROL_WORD_ADDRESS,
		data.PUMP_3_STATUS_WORD_ADDRESS)
	stopActions <- pump3StopAction

	pump4StopAction := action.NewStopAction(data.PUMP_4_CONTROL_WORD_ADDRESS,
		data.PUMP_4_STATUS_WORD_ADDRESS)
	stopActions <- pump4StopAction

	pump5StopAction := action.NewStopAction(data.PUMP_5_CONTROL_WORD_ADDRESS,
		data.PUMP_5_STATUS_WORD_ADDRESS)
	stopActions <- pump5StopAction

	pump6StopAction := action.NewStopAction(data.PUMP_6_CONTROL_WORD_ADDRESS,
		data.PUMP_6_STATUS_WORD_ADDRESS)
	stopActions <- pump6StopAction

	pump7StopAction := action.NewStopAction(data.PUMP_7_CONTROL_WORD_ADDRESS,
		data.PUMP_7_STATUS_WORD_ADDRESS)
	stopActions <- pump7StopAction

	pump8StopAction := action.NewStopAction(data.PUMP_8_CONTROL_WORD_ADDRESS,
		data.PUMP_8_STATUS_WORD_ADDRESS)
	stopActions <- pump8StopAction

	pump9StopAction := action.NewStopAction(data.PUMP_9_CONTROL_WORD_ADDRESS,
		data.PUMP_9_STATUS_WORD_ADDRESS)
	stopActions <- pump9StopAction

	pump10StopAction := action.NewStopAction(data.PUMP_10_CONTROL_WORD_ADDRESS,
		data.PUMP_10_STATUS_WORD_ADDRESS)
	stopActions <- pump10StopAction

	close(stopActions)

	//go func() {
	//	for i := 0; i < len(stopActions); i++ {
	//		controller.ExecuteImmediately(<-stopActions)
	//		if i == 4 {
	//			time.Sleep(200 * time.Millisecond)
	//		}
	//	}
	//}()

	go func() {
		for stopAction := range stopActions {
			time.Sleep(100 * time.Millisecond)
			controller.ExecuteImmediately(stopAction)
		}
	}()

	//go func() {
	//	controller.ExecuteImmediately(axisR1StopAction)
	//	controller.ExecuteImmediately(axisXStopAction)
	//	controller.ExecuteImmediately(axisYStopAction)
	//	controller.ExecuteImmediately(axisZStopAction)
	//	controller.ExecuteImmediately(temperatureResetAction)
	//	controller.ExecuteImmediately(dishOutStopAction)
	//	controller.ExecuteImmediately(shakeStopAction)
	//	controller.ExecuteImmediately(pump1StopAction)
	//	controller.ExecuteImmediately(pump2StopAction)
	//	controller.ExecuteImmediately(pump3StopAction)
	//	controller.ExecuteImmediately(pump4StopAction)
	//	controller.ExecuteImmediately(pump5StopAction)
	//	controller.ExecuteImmediately(pump6StopAction)
	//	controller.ExecuteImmediately(pump7StopAction)
	//	controller.ExecuteImmediately(pump8StopAction)
	//	controller.ExecuteImmediately(pump9StopAction)
	//	controller.ExecuteImmediately(pump10StopAction)
	//}()

}
