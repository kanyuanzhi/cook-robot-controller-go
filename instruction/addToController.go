package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/config"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (ins Instruction) AddToController(controller *core.Controller) {
	return
}

func (ins IngredientInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_INGREDIENT_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	num, err := strconv.ParseUint(ins.SlotNumber, 10, 32)
	if err != nil {
		fmt.Println("无法将字符串转换为uint32")
		return
	}
	slotNumber := uint32(num)
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.SlotNumberToPosition[slotNumber]]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, data.SHAKE_AMOUNT),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, data.SHAKE_UPWARD_SPEED),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, data.SHAKE_DOWNWARD_SPEED),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, data.SHAKE_DISTANCE))
	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(shakeControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 6))
	logger.Log.Printf("[步骤]添加食材，打开%d号菜仓。", slotNumber)
}

func (ins SeasoningInstruction) AddToController(controller *core.Controller) {
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
	for pump, weight := range ins.PumpToWeightMap {
		num, err := strconv.ParseUint(pump, 10, 32)
		if err != nil {
			logger.Log.Println("无法将字符串转换为uint32")
			return
		}
		pumpNumber := uint32(num)
		duration := weight * ins.PumpToRatioMap[pump]

		//var duration uint32 = 0
		//if pumpNumber == 1 { // 液体泵（食用油）
		//	duration = uint32(math.Round(float64(weight) * 1000 / 12.8))
		//} else if pumpNumber >= 2 && pumpNumber <= 5 { // 液体泵（其他液体调料）
		//	duration = uint32(math.Round(float64(weight) * 1000 / 6.4))
		//} else if pumpNumber == 6 { // 液体泵（纯净水）
		//	duration = uint32(math.Round(float64(weight) * 1000 / 12.8))
		//} else if pumpNumber >= 7 && pumpNumber <= 8 { // 水泵，暂不用
		//	duration = uint32(math.Round(float64(weight) * 1000 / 6.4))
		//} else if pumpNumber == 9 { // 固体泵（食盐）
		//	duration = weight * 1000 / 10
		//} else if pumpNumber == 10 { // 固体泵（鸡精）
		//	duration = weight * 1000 / 5
		//} else {
		//	logger.Log.Println("wrong pump number")
		//	return
		//}
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
	actionNumber := 0
	if len(liquidGroupAction.Actions) != 0 {
		controller.AddAction(axisYLocateControlActionLiquid)
		controller.AddAction(liquidGroupAction)
		actionNumber += 3
	}
	if len(solidGroupAction.Actions) != 0 {
		controller.AddAction(axisYLocateControlActionSolid)
		controller.AddAction(solidGroupAction)
		actionNumber += 3
	}
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, actionNumber))
	logger.Log.Printf("[步骤]添加调料%s", logStr)
}

func (ins WaterInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	duration := ins.Weight * ins.Ratio
	//var duration uint32 = 0
	//if ins.PumpNumber >= 6 && ins.PumpNumber <= 8 { // 水泵
	//	duration = uint32(math.Round(float64(ins.Weight) * 1000 / 12.8))
	//} else {
	//	logger.Log.Println("wrong pump number")
	//}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[ins.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[ins.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[ins.PumpNumber], duration))

	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddAction(pumpControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 4))
	logger.Log.Printf("[步骤]添加水，%d号泵打开%d毫秒（%d克）", ins.PumpNumber, duration, ins.Weight)
}

func (ins OilInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	duration := ins.Weight * ins.Ratio
	//var duration uint32 = 0
	//if ins.PumpNumber == 1 { // 油泵
	//	duration = uint32(math.Round(float64(ins.Weight) * 1000 / 12.8))
	//} else {
	//	logger.Log.Println("wrong pump number")
	//}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[ins.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[ins.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[ins.PumpNumber], duration))

	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddAction(pumpControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 4))
	logger.Log.Printf("[步骤]添加油，%d号泵打开%d毫秒（%d克）", ins.PumpNumber, duration, ins.Weight)
}

func (ins StirFryInstruction) AddToController(controller *core.Controller) {
	speed := ins.Gear * data.R1_MAX_ROTATE_SPEED / 5
	axleRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, speed),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	delayAction := action.NewDelayAction(ins.Duration * 1000)

	// 先旋转再定位y轴
	controller.AddAction(axleRotateControlAction)
	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(delayAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 5))
	logger.Log.Printf("[步骤]添加翻炒，%d号档位（转速%d）持续%d秒", ins.Gear, speed, ins.Duration)
}

func (ins HeatInstruction) AddToController(controller *core.Controller) {
	temperatureControlAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, uint32(math.Round(ins.Temperature*10))))
	controller.AddAction(temperatureControlAction)
	var judgeTypeStr string
	switch ins.JudgeType {
	case BOTTOM_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(
			data.NewAddressValue(data.TEMPERATURE_BOTTOM_ADDRESS, uint32(math.Round(ins.TargetTemperature*10))), data.LARGER_THAN_TARGET, data.TEMPERATURE_CONTROL_WORD_ADDRESS)
		controller.AddAction(triggerAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
		judgeTypeStr = fmt.Sprintf("锅底温度达到%.1f℃后开始下一步骤", ins.TargetTemperature)
		break
	case INFRARED_TEMPERATURE_JUDGE_TYPE:
		axisYLocateStirFryControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
			data.Y_LOCATE_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
			data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
		triggerAction := action.NewTriggerAction(
			data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, uint32(math.Round(ins.TargetTemperature*10))), data.LARGER_THAN_TARGET, data.TEMPERATURE_CONTROL_WORD_ADDRESS)
		controller.AddAction(axisYLocateStirFryControlAction) // 红外测温前确保y轴前往翻炒位
		controller.AddAction(triggerAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 5))
		judgeTypeStr = fmt.Sprintf("红外温度达到%.1f℃后开始下一步骤", ins.TargetTemperature)
		break
	case DURATION_JUDGE_TYPE:
		delayAction := action.NewDelayAction(ins.Duration * 1000)
		controller.AddAction(delayAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
		judgeTypeStr = fmt.Sprintf("加热%d秒开始下一步骤", ins.Duration)
		break
	case NO_JUDGE:
		judgeTypeStr = "不判断"
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
		break
	default:
		delayAction := action.NewDelayAction(1000)
		controller.AddAction(delayAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
		logger.Log.Println("wrong temperature judge type")
	}
	logger.Log.Printf("[步骤]添加火力，加热温度%.1f℃，%s", ins.Temperature, judgeTypeStr)
}

func (ins DishOutInstruction) AddToController(controller *core.Controller) {
	dishOutControlAction := action.NewDishOutControlAction(data.DISH_OUT_CONTROL_WORD_ADDRESS,
		data.DISH_OUT_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.DISH_OUT_AMOUNT_ADDRESS, data.DISH_OUT_AMOUNT),
		data.NewAddressValue(data.DISH_OUT_UPWARD_SPEED_ADDRESS, data.DISH_OUT_UPWARD_SPEED),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_SPEED_ADDRESS, data.DISH_OUT_DOWNWARD_SPEED),
		data.NewAddressValue(data.DISH_OUT_UPWARD_POSITION_ADDRESS, data.YPositionToDistance[data.Y_DISH_OUT_HIGH_POSITION]),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_POSITION_ADDRESS, data.YPositionToDistance[data.Y_DISH_OUT_LOW_POSITION]))
	controller.AddAction(dishOutControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加出菜，次数%d，上行速度%d，下行速度%d，上行位置%d，下行位置%d",
		data.DISH_OUT_AMOUNT, data.DISH_OUT_UPWARD_SPEED, data.DISH_OUT_DOWNWARD_SPEED,
		data.YPositionToDistance[data.Y_DISH_OUT_HIGH_POSITION], data.YPositionToDistance[data.Y_DISH_OUT_LOW_POSITION])
}

func (ins ShakeInstruction) AddToController(controller *core.Controller) {
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, data.SHAKE_AMOUNT),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, data.SHAKE_UPWARD_SPEED),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, data.SHAKE_DOWNWARD_SPEED),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, data.SHAKE_DISTANCE))
	controller.AddAction(shakeControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加抖菜，次数%d，上行速度%d，下行速度%d，抖动距离%d",
		data.SHAKE_AMOUNT, data.SHAKE_UPWARD_SPEED, data.SHAKE_DOWNWARD_SPEED, data.SHAKE_DISTANCE)
}

func (ins LampblackPurifyInstruction) AddToController(controller *core.Controller) {
	lampblackPurifyControlAction := action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
		data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, ins.Mode))
	controller.AddAction(lampblackPurifyControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加油烟净化，模式%d", ins.Mode)
}

func (ins DoorUnlockInstruction) AddToController(controller *core.Controller) {
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.AddAction(doorUnlockControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加打开电磁门锁")
}

func (ins InitInstruction) AddToController(controller *core.Controller) {
	rotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, data.R1_MAX_ROTATE_SPEED/5), // 炒菜启动前R1轴自动以1档转动
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	var lampblackPurifyControlAction action.Actioner
	if config.Parameter.LampblackPurify.Enable {
		lampblackPurifyControlAction = action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
			data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, data.LAMPBLACK_PURIFY_PURIFICATION_MODE))
	} else {
		lampblackPurifyControlAction = action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
			data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, data.LAMPBLACK_PURIFY_VENTING_MODE))
	}

	delayAction := action.NewDelayAction(2000)
	actionNumber := 3
	controller.AddAction(rotateControlAction)
	if data.LAMPBLACK_PURIFY_ENABLE {
		controller.AddAction(lampblackPurifyControlAction)
		actionNumber += 2
	}
	controller.AddAction(delayAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, actionNumber))
	logger.Log.Printf("[步骤]添加转动，油烟净化设置为模式%d（视机器支持情况）", data.LAMPBLACK_PURIFY_PURIFICATION_MODE)
}

func (ins FinishInstruction) AddToController(controller *core.Controller) {
	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	axisRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 100),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_READY_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	delayAction := action.NewDelayAction(2000)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)

	var lampblackPurifyControlAction action.Actioner
	if data.LAMPBLACK_PURIFY_AUTO_START {
		lampblackPurifyControlAction = action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
			data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, data.LAMPBLACK_PURIFY_VENTING_MODE))
	} else {
		lampblackPurifyControlAction = action.NewStopAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
			data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS)
	}

	actionNumber := 11
	controller.AddAction(temperatureResetAction)
	controller.AddAction(axisRotateControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(axisR1StopAction)
	controller.AddAction(temperatureResetAction) // 最后再次添加一次温控停止，确保温控停止
	if data.LAMPBLACK_PURIFY_ENABLE {
		controller.AddAction(lampblackPurifyControlAction)
		actionNumber += 2
	}
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, actionNumber))
	logger.Log.Printf("[步骤]添加转动、温控停止，X轴回备菜位，油烟净化设置为模式%d（视机器支持情况）", data.LAMPBLACK_PURIFY_VENTING_MODE)
}

func (ins ResetXYTInstruction) AddToController(controller *core.Controller) {
	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	axisXResetControlAction := action.NewAxisResetControlAction(data.X_RESET_CONTROL_WORD_ADDRESS,
		data.X_RESET_STATUS_WORD_ADDRESS)
	axisYResetControlAction := action.NewAxisResetControlAction(data.Y_RESET_CONTROL_WORD_ADDRESS,
		data.Y_RESET_STATUS_WORD_ADDRESS)

	controller.AddAction(temperatureResetAction)
	controller.AddAction(axisXResetControlAction)
	controller.AddAction(axisYResetControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 6))
	logger.Log.Printf("[步骤]添加X轴、Y轴复位，温控停止")
}

func (ins ResetRTInstruction) AddToController(controller *core.Controller) {
	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	axisRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 100),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_READY_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	delayAction := action.NewDelayAction(2000)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)

	controller.AddAction(temperatureResetAction)
	controller.AddAction(axisRotateControlAction)
	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(axisR1StopAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 9))
	logger.Log.Printf("[步骤]添加转动、温控停止，X轴回备菜位")
}

func (ins PrepareInstruction) AddToController(controller *core.Controller) {
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_READY_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(axisYLocateControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 4))
	logger.Log.Printf("[步骤]添加X轴、Y轴定位到备菜位")
}

func (ins WashInstruction) AddToController(controller *core.Controller) {
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_WITHDRAWER_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_WASH_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	axisRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 3),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, data.WASH_ROTATE_SPEED),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, data.WASH_ROTATE_CROSS_AMOUNT))
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[data.WASH_PUMP_NUMBER],
		data.PumpNumberToPumpStatusWordAddress[data.WASH_PUMP_NUMBER],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[data.WASH_PUMP_NUMBER], data.WASH_ADD_WATER_DURATION))
	//temperatureControlAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
	//	data.TEMPERATURE_STATUS_WORD_ADDRESS,
	//	data.NewAddressValue(data.TEMPERATURE_ADDRESS, data.WASH_TEMPERATURE))
	delayAction := action.NewDelayAction(data.WASH_DURATION)
	temperatureResetAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, 0))
	slowAxisRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, 100),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	slowDelayAction := action.NewDelayAction(2000)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)
	axisYLocateStirFryControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	controller.AddAction(axisXLocateControlAction)
	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axisRotateControlAction)
	controller.AddAction(pumpControlAction)
	//controller.AddAction(temperatureControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(temperatureResetAction)
	controller.AddAction(slowAxisRotateControlAction)
	controller.AddAction(slowDelayAction)
	controller.AddAction(axisR1StopAction)
	controller.AddAction(axisYLocateStirFryControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 18))
	logger.Log.Printf("[步骤]添加洗锅")
}

func (ins PourInstruction) AddToController(controller *core.Controller) {
	axisYLocateWashControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_POUR_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	delayAction := action.NewDelayAction(data.WASH_POUR_WATER_DURATION)
	axisYLocateStirFryControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))

	controller.AddAction(axisYLocateWashControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(axisYLocateStirFryControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 5))
	logger.Log.Printf("[步骤]添加倒水")
}

func (ins DelayInstruction) AddToController(controller *core.Controller) {
	delayAction := action.NewDelayAction(ins.Duration)
	controller.AddAction(delayAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 1))
	logger.Log.Printf("[步骤]延时%d秒")
}

func (ins WithdrawInstruction) AddToController(controller *core.Controller) {
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[data.X_WITHDRAWER_POSITION]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))

	controller.AddAction(axisXLocateControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加X轴定位到收纳位")
}

// straight instruction

func (ins AxisInstruction) AddToController(controller *core.Controller) {
	if ins.Function != RESET && ins.Function != LOCATE {
		logger.Log.Println("wrong axis function")
		return
	}
	ins.Axis = strings.ToUpper(ins.Axis)
	if ins.Axis != "X" && ins.Axis != "Y" && ins.Axis != "Z" && ins.Axis != "R1" && ins.Axis != "R2" {
		logger.Log.Println("wrong axis")
		return
	}
	var axisControlAction action.Actioner
	if ins.Function == RESET { // 复位
		axisControlAction = action.NewAxisResetControlAction(data.AxisToResetControlWordAddress[ins.Axis],
			data.AxisToResetStatusWordAddress[ins.Axis])
	} else { // 定位
		axisControlAction = action.NewAxisLocateControlAction(data.AxisToLocateControlWordAddress[ins.Axis],
			data.AxisToLocateStatusWordAddress[ins.Axis],
			data.NewAddressValue(data.AxisToLocatePositionAddress[ins.Axis], ins.TargetPosition),
			data.NewAddressValue(data.AxisToLocateSpeedAddress[ins.Axis], ins.Speed))
	}
	controller.AddAction(axisControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
}

func (ins RotateInstruction) AddToController(controller *core.Controller) {
	if ins.Function != STOP && ins.Function != START {
		logger.Log.Println("wrong rotate function")
		return
	}
	var rotateControlAction action.Actioner
	if ins.Function == STOP { // 停转
		rotateControlAction = action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
			data.R1_ROTATE_STATUS_WORD_ADDRESS)
		controller.AddAction(rotateControlAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	} else { // 旋转
		rotateControlAction = action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
			data.R1_ROTATE_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, ins.Mode),
			data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, ins.Speed),
			data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, ins.RotationalAmount))
		delayAction := action.NewDelayAction(2000)
		controller.AddAction(rotateControlAction)
		controller.AddAction(delayAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
	}
	if ins.Function == STOP {
		logger.Log.Printf("[步骤]添加R轴停转")
	} else {
		logger.Log.Printf("[步骤]添加R轴%s，速度%d，正反转圈数%d", data.RotateModeToString[ins.Mode], ins.Speed, ins.RotationalAmount)
	}
}
