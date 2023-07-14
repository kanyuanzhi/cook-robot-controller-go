package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"fmt"
	"math"
	"strconv"
)

func (i Instruction) AddToController(controller *core.Controller) {
	return
}

func (i IngredientInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_INGREDIENT_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	num, err := strconv.ParseUint(i.SlotNumber, 10, 32)
	if err != nil {
		fmt.Println("无法将字符串转换为uint32")
		return
	}
	slotNumber := uint32(num)
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[slotNumber]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 2000))
	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(shakeControlAction)
	logger.Log.Printf("[步骤]添加食材，打开%d号菜仓。", slotNumber)
}

func (s SeasoningInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	axisYLocateControlActionSolid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_SOLID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	liquidGroupAction := action.NewGroupAction()
	solidGroupAction := action.NewGroupAction()
	var logStr string
	for pumpNumber, weight := range s.PumpToWeightMap {
		num, err := strconv.ParseUint(pumpNumber, 10, 32)
		if err != nil {
			logger.Log.Println("无法将字符串转换为uint32")
			return
		}
		pumpNumber := uint32(num)
		var duration uint32 = 0
		if pumpNumber >= 1 && pumpNumber <= 6 { // 液体泵
			duration = uint32(math.Round(float64(weight) * 1000 / 6.4))
		} else if pumpNumber >= 7 && pumpNumber <= 8 { // 水泵
			duration = uint32(math.Round(float64(weight) * 1000 / 6.4))
		} else if pumpNumber >= 9 && pumpNumber <= 10 { // 固体泵
			duration = weight * 100
		} else {
			logger.Log.Println("wrong pump number")
			return
		}
		logStr += fmt.Sprintf("，%d号泵打开%d毫秒（%d克）", pumpNumber, duration, weight)
		pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[pumpNumber],
			data.PumpNumberToPumpStatusWordAddress[pumpNumber],
			data.NewAddressValue(data.PumpNumberToPumpDurationAddress[pumpNumber], duration))
		if pumpNumber >= 1 && pumpNumber <= 8 { // 液体泵和水泵
			liquidGroupAction.AddAction(pumpControlAction)
		} else { // 固体泵
			solidGroupAction.AddAction(pumpControlAction)
		}
	}
	if len(liquidGroupAction.Actions) != 0 {
		controller.AddAction(axisYLocateControlActionLiquid)
		controller.AddAction(liquidGroupAction)
	}
	if len(solidGroupAction.Actions) != 0 {
		controller.AddAction(axisYLocateControlActionSolid)
		controller.AddAction(solidGroupAction)
	}
	logger.Log.Printf("[步骤]添加调料%s", logStr)
}

func (w WaterInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	var duration uint32 = 0
	if w.PumpNumber >= 7 && w.PumpNumber <= 8 { // 水泵
		duration = uint32(math.Round(float64(w.Weight) * 1000 / 6.4))
	} else {
		logger.Log.Println("wrong pump number")
	}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[w.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[w.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[w.PumpNumber], duration))

	controller.AddAction(pumpControlAction)
	controller.AddAction(axisYLocateControlActionLiquid)
	logger.Log.Printf("[步骤]添加水，%d号泵打开%d毫秒（%d克）", w.PumpNumber, duration, w.Weight)
}

func (o OilInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	var duration uint32 = 0
	if o.PumpNumber == 1 { // 油泵
		duration = uint32(math.Round(float64(o.Weight) * 1000 / 6.4))
	} else {
		logger.Log.Println("wrong pump number")
	}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[o.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[o.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[o.PumpNumber], duration))

	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddAction(pumpControlAction)
	logger.Log.Printf("[步骤]添加油，%d号泵打开%d毫秒（%d克）", o.PumpNumber, duration, o.Weight)
}

func (s StirFryInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	speedValue := s.Gear * 350
	axleRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, speedValue),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	delayAction := action.NewDelayAction(s.Duration * 1000)

	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axleRotateControlAction)
	controller.AddAction(delayAction)
	logger.Log.Printf("[步骤]添加翻炒，%d号档位（转速%d）持续%d秒", s.Gear, speedValue, s.Duration)
}

func (h HeatInstruction) AddToController(controller *core.Controller) {
	temperatureControlAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, uint32(math.Round(h.Temperature*10))))
	controller.AddAction(temperatureControlAction)
	var judgeTypeStr string
	switch h.JudgeType {
	case BOTTOM_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(
			data.NewAddressValue(data.TEMPERATURE_BOTTOM_ADDRESS, uint32(math.Round(h.TargetTemperature*10))), data.LARGER_THAN_TARGET)
		controller.AddAction(triggerAction)
		judgeTypeStr = fmt.Sprintf("锅底温度达到%.1f℃后开始下一步骤", h.TargetTemperature)
		break
	case INFRARED_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(
			data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, uint32(math.Round(h.TargetTemperature*10))), data.EQUAL_TO_TARGET)
		controller.AddAction(triggerAction)
		judgeTypeStr = fmt.Sprintf("红外温度达到%.1f℃后开始下一步骤", h.TargetTemperature)
		break
	case DURATION_JUDGE_TYPE:
		delayAction := action.NewDelayAction(h.Duration * 1000)
		controller.AddAction(delayAction)
		judgeTypeStr = fmt.Sprintf("加热%d秒开始下一步骤", h.Duration)
		break
	case NO_JUDGE:
		judgeTypeStr = "不判断"
		break
	default:
		delayAction := action.NewDelayAction(1000)
		controller.AddAction(delayAction)
		logger.Log.Println("wrong temperature judge type")
	}
	logger.Log.Printf("[步骤]添加火力，加热温度%.1f℃，%s", h.Temperature, judgeTypeStr)
}

func (d DishOutInstruction) AddToController(controller *core.Controller) {
	dishOutControlAction := action.NewDishOutControlAction(data.DISH_OUT_CONTROL_WORD_ADDRESS,
		data.DISH_OUT_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.DISH_OUT_AMOUNT_ADDRESS, 3),
		data.NewAddressValue(data.DISH_OUT_UPWARD_SPEED_ADDRESS, 400),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_SPEED_ADDRESS, 800),
		data.NewAddressValue(data.DISH_OUT_UPWARD_POSITION_ADDRESS, 2000),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_POSITION_ADDRESS, 4000))
	controller.AddAction(dishOutControlAction)
	logger.Log.Printf("[步骤]添加出菜，次数%d，上行速度%d，下行速度%d，上行位置%d，下行位置%d", 3, 400, 800, 2000, 4000)
}

func (s ShakeInstruction) AddToController(controller *core.Controller) {
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 2000))
	controller.AddAction(shakeControlAction)
	logger.Log.Printf("[步骤]添加抖菜，次数%d，上行速度%d，下行速度%d，抖动距离%d", 5, 30000, 20000, 2000)
}

func (l LampblackPurifyInstruction) AddToController(controller *core.Controller) {
	lampblackPurifyControlAction := action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
		data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, l.Mode))
	controller.AddAction(lampblackPurifyControlAction)
	logger.Log.Printf("[步骤]添加油烟净化，模式%d", l.Mode)
}

func (d DoorUnlockInstruction) AddToController(controller *core.Controller) {
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.AddAction(doorUnlockControlAction)
	logger.Log.Printf("[步骤]添加打开电磁门锁")
}

func (r ResetAllInstruction) AddToController(controller *core.Controller) {
	//temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
	//	data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	//axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
	//	data.R1_ROTATE_STATUS_WORD_ADDRESS)
	axisXResetControlAction := action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS,
		data.X_RESET_STATUS_WORD_ADDRESS)
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_READY_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	axisYResetControlAction := action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS,
		data.Y_RESET_STATUS_WORD_ADDRESS)
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_2_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	//controller.AddAction(temperatureResetAction)
	//controller.AddAction(axisR1StopAction)
	controller.AddAction(axisXResetControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(axisYResetControlAction)
	controller.AddAction(axisYLocateControlAction)
}

func (r ResetXYInstruction) AddToController(controller *core.Controller) {
	axisXResetControlAction := action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS,
		data.X_RESET_STATUS_WORD_ADDRESS)
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_READY_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	axisYResetControlAction := action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS,
		data.Y_RESET_STATUS_WORD_ADDRESS)
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_2_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	controller.AddAction(axisXResetControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(axisYResetControlAction)
	controller.AddAction(axisYLocateControlAction)
	logger.Log.Printf("[步骤]添加X轴、Y轴复位")
}

func (r ResetRTInstruction) AddToController(controller *core.Controller) {
	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	axleRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 100),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	delayAction := action.NewDelayAction(1000)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)

	controller.AddAction(temperatureResetAction)
	controller.AddAction(axleRotateControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(axisR1StopAction)
	logger.Log.Printf("[步骤]添加转动、温控停止")
}
