package instruction

import (
	"cook-robot-controller-go/action"
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
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[slotNumber]),
		data.NewAddressValue(data.X_LOCATE_SPEED_ADDRESS, data.X_MOVE_SPEED))
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 3200))
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
	for pumpNumber, weight := range ins.PumpToWeightMap {
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
	var duration uint32 = 0
	if ins.PumpNumber >= 7 && ins.PumpNumber <= 8 { // 水泵
		duration = uint32(math.Round(float64(ins.Weight) * 1000 / 6.4))
	} else {
		logger.Log.Println("wrong pump number")
	}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[ins.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[ins.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[ins.PumpNumber], duration))

	controller.AddAction(pumpControlAction)
	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 4))
	logger.Log.Printf("[步骤]添加水，%d号泵打开%d毫秒（%d克）", ins.PumpNumber, duration, ins.Weight)
}

func (ins OilInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_LIQUID_SEASONING_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	var duration uint32 = 0
	if ins.PumpNumber == 1 { // 油泵
		duration = uint32(math.Round(float64(ins.Weight) * 1000 / 6.4))
	} else {
		logger.Log.Println("wrong pump number")
	}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[ins.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[ins.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[ins.PumpNumber], duration))

	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddAction(pumpControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 4))
	logger.Log.Printf("[步骤]添加油，%d号泵打开%d毫秒（%d克）", ins.PumpNumber, duration, ins.Weight)
}

func (ins StirFryInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_STIR_FRY_3_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	speedValue := ins.Gear * 350
	axleRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, speedValue),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	delayAction := action.NewDelayAction(ins.Duration * 1000)

	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axleRotateControlAction)
	controller.AddAction(delayAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 5))
	logger.Log.Printf("[步骤]添加翻炒，%d号档位（转速%d）持续%d秒", ins.Gear, speedValue, ins.Duration)
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
			data.NewAddressValue(data.TEMPERATURE_BOTTOM_ADDRESS, uint32(math.Round(ins.TargetTemperature*10))), data.LARGER_THAN_TARGET, data.TEMPERATURE_CONTROL_WORD_ADDRESS, 0)
		controller.AddAction(triggerAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
		judgeTypeStr = fmt.Sprintf("锅底温度达到%.1f℃后开始下一步骤", ins.TargetTemperature)
		break
	case INFRARED_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(
			data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, uint32(math.Round(ins.TargetTemperature*10))), data.LARGER_THAN_TARGET, data.TEMPERATURE_CONTROL_WORD_ADDRESS, 0)
		controller.AddAction(triggerAction)
		controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 3))
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
		data.NewAddressValue(data.DISH_OUT_AMOUNT_ADDRESS, 3),
		data.NewAddressValue(data.DISH_OUT_UPWARD_SPEED_ADDRESS, 600),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_SPEED_ADDRESS, 2200),
		data.NewAddressValue(data.DISH_OUT_UPWARD_POSITION_ADDRESS, data.YPositionToDistance[data.Y_DISH_OUT_HEIGH_POSITION]),
		data.NewAddressValue(data.DISH_OUT_DOWNWARD_POSITION_ADDRESS, data.YPositionToDistance[data.Y_DISH_OUT_LOW_POSITION]))
	controller.AddAction(dishOutControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加出菜，次数%d，上行速度%d，下行速度%d，上行位置%d，下行位置%d", 3, 400, 800,
		data.YPositionToDistance[data.Y_DISH_OUT_HEIGH_POSITION], data.YPositionToDistance[data.Y_DISH_OUT_LOW_POSITION])
}

func (ins ShakeInstruction) AddToController(controller *core.Controller) {
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 3200))
	controller.AddAction(shakeControlAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 2))
	logger.Log.Printf("[步骤]添加抖菜，次数%d，上行速度%d，下行速度%d，抖动距离%d", 5, 30000, 20000, 3200)
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
	delayAction := action.NewDelayAction(2000)
	axisR1StopAction := action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS)

	controller.AddAction(temperatureResetAction)
	controller.AddAction(axisRotateControlAction)
	controller.AddAction(delayAction)
	controller.AddAction(axisR1StopAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 7))
	logger.Log.Printf("[步骤]添加转动、温控停止")
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

func (ins DelayInstruction) AddToController(controller *core.Controller) {
	delayAction := action.NewDelayAction(ins.Duration)
	controller.AddAction(delayAction)
	controller.AddInstructionInfo(data.NewInstructionInfo(string(ins.InstructionType), ins.InstructionName, 1))
	logger.Log.Printf("[步骤]延时%d秒")
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
