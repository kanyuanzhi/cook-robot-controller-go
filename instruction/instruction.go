package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"math"
)

type Instructioner interface {
	CheckType() InstructionType
	AddToController(controller *core.Controller)
}

type InstructionType string

const (
	INGREDIENT       = InstructionType("ingredient")
	SEASONING        = InstructionType("seasoning")
	WATER            = InstructionType("water")
	STIR_FRY         = InstructionType("stir_fry")
	HEAT             = InstructionType("heat")
	DISH_OUT         = InstructionType("dish_out")
	SHAKE            = InstructionType("shake")
	LAMPBLACK_PURIFY = InstructionType("lampblack_purify")
	DOOR_UNLOCK      = InstructionType("door_unlock")

	AXIS   = InstructionType("axis")
	ROTATE = InstructionType("rotate")
	PUMP   = InstructionType("pump")
)

type Instruction struct {
	InstructionType InstructionType `json:"instruction_type" mapstructure:"instruction_type"`
}

func (i Instruction) CheckType() InstructionType {
	return i.InstructionType
}

func (i Instruction) AddToController(controller *core.Controller) {
	return
}

func NewInstruction(instructionType InstructionType) Instruction {
	return Instruction{InstructionType: instructionType}
}

type IngredientInstruction struct {
	Instruction `mapstructure:",squash"`
	SlotNumber  uint32 `json:"slot_number" mapstructure:"slot_number"`
}

func NewIngredientInstruction(slotNumber uint32) *IngredientInstruction {
	return &IngredientInstruction{
		Instruction: NewInstruction(INGREDIENT),
		SlotNumber:  slotNumber,
	}
}

func (i IngredientInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.YPositionToDistance[data.Y_INGREDIENT_POSITION]),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	axisXLocateControlAction := action.NewAxisLocateControlAction(data.X_LOCATE_CONTROL_WORD_ADDRESS,
		data.X_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.X_LOCATE_POSITION_ADDRESS, data.XPositionToDistance[i.SlotNumber]),
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
}

type SeasoningInstruction struct {
	Instruction     `mapstructure:",squash"`
	PumpToWeightMap map[uint32]uint32 `json:"pump_to_weight_map" mapstructure:"pump_to_weight_map"` // 泵号:重量
}

func NewSeasoningInstruction(pumpToWeightMap map[uint32]uint32) *SeasoningInstruction {
	return &SeasoningInstruction{
		Instruction:     NewInstruction(SEASONING),
		PumpToWeightMap: pumpToWeightMap,
	}
}

func (s SeasoningInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlActionLiquid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.Y_LIQUID_SEASONING_POSITION),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	axisYLocateControlActionSolid := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.Y_SOLID_SEASONING_POSITION),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	liquidGroupAction := action.NewGroupAction()
	solidGroupAction := action.NewGroupAction()
	for pumpNumber, weight := range s.PumpToWeightMap {
		var duration uint32 = 0
		if pumpNumber >= 1 && pumpNumber <= 6 { // 液体泵
			duration = weight * 10
		} else if pumpNumber >= 7 && pumpNumber <= 8 { // 水泵
			duration = weight * 100
		} else if pumpNumber >= 9 && pumpNumber <= 10 { // 固体泵
			duration = weight * 100
		} else {
			logger.Log.Println("wrong pump number")
			return
		}
		pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[pumpNumber],
			data.PumpNumberToPumpStatusWordAddress[pumpNumber],
			data.NewAddressValue(data.PumpNumberToPumpDurationAddress[pumpNumber], duration))
		if pumpNumber >= 1 && pumpNumber <= 8 { // 液体泵和水泵
			liquidGroupAction.AddAction(pumpControlAction)
		} else { // 固体泵
			solidGroupAction.AddAction(pumpControlAction)
			return
		}
	}

	controller.AddAction(axisYLocateControlActionLiquid)
	controller.AddAction(liquidGroupAction)
	controller.AddAction(axisYLocateControlActionSolid)
	controller.AddAction(solidGroupAction)
}

type WaterInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pump_number" mapstructure:"pump_number"`
	Weight      uint32 `json:"weight"`
}

func NewWaterInstruction(pumpNumber uint32, weight uint32) *WaterInstruction {
	return &WaterInstruction{
		Instruction: NewInstruction(WATER),
		PumpNumber:  pumpNumber,
		Weight:      weight,
	}
}

func (w WaterInstruction) AddToController(controller *core.Controller) {
	var duration uint32 = 0
	if w.PumpNumber >= 7 && w.PumpNumber <= 8 { // 液体泵
		duration = w.Weight * 10
	} else {
		logger.Log.Println("wrong pump number")
	}
	pumpControlAction := action.NewPumpControlAction(data.PumpNumberToPumpControlWordAddress[w.PumpNumber],
		data.PumpNumberToPumpStatusWordAddress[w.PumpNumber],
		data.NewAddressValue(data.PumpNumberToPumpDurationAddress[w.PumpNumber], duration))

	controller.AddAction(pumpControlAction)
}

type StirFryInstruction struct {
	Instruction `mapstructure:",squash"`
	Gear        uint32 `json:"gear"`
	Duration    uint32 `json:"duration"`
}

func NewStirFryInstruction(gear uint32, duration uint32) *StirFryInstruction {
	return &StirFryInstruction{
		Instruction: NewInstruction(STIR_FRY),
		Gear:        gear,
		Duration:    duration,
	}
}

func (s StirFryInstruction) AddToController(controller *core.Controller) {
	axisYLocateControlAction := action.NewAxisLocateControlAction(data.Y_LOCATE_CONTROL_WORD_ADDRESS,
		data.Y_LOCATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.Y_LOCATE_POSITION_ADDRESS, data.Y_STIR_FRY_1_POSITION),
		data.NewAddressValue(data.Y_LOCATE_SPEED_ADDRESS, data.Y_MOVE_SPEED))
	speedValue := s.Gear * 200
	axleRotateControlAction := action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
		data.R1_ROTATE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, 1),
		data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, speedValue),
		data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, 0))
	delayAction := action.NewDelayAction(s.Duration * 1000)

	controller.AddAction(axisYLocateControlAction)
	controller.AddAction(axleRotateControlAction)
	controller.AddAction(delayAction)
}

type HeatingInstruction struct {
	Instruction       `mapstructure:",squash"`
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"target_temperature" mapstructure:"target_temperature"`
	Duration          uint32  `json:"duration"`
	JudgeType         uint    `json:"judge_type" mapstructure:"judge_type"`
}

func NewHeatingInstruction(temperature float64, targetTemperature float64, duration uint32, judgeType uint) *HeatingInstruction {
	return &HeatingInstruction{
		Instruction:       NewInstruction(HEAT),
		Temperature:       temperature,
		TargetTemperature: targetTemperature,
		Duration:          duration,
		JudgeType:         judgeType,
	}
}

const (
	BOTTOM_TEMPERATURE_JUDGE_TYPE uint = iota + 1
	INFRARED_TEMPERATURE_JUDGE_TYPE
	DURATION_JUDGE_TYPE
)

func (h HeatingInstruction) AddToController(controller *core.Controller) {
	temperatureControlAction := action.NewTemperatureControlAction(data.TEMPERATURE_CONTROL_WORD_ADDRESS,
		data.TEMPERATURE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.TEMPERATURE_ADDRESS, uint32(math.Round(h.Temperature*10))))
	controller.AddAction(temperatureControlAction)
	switch h.JudgeType {
	case BOTTOM_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(data.NewAddressValue(data.TEMPERATURE_BOTTOM_ADDRESS, uint32(math.Round(h.TargetTemperature*10))))
		controller.AddAction(triggerAction)
		break
	case INFRARED_TEMPERATURE_JUDGE_TYPE:
		triggerAction := action.NewTriggerAction(data.NewAddressValue(data.TEMPERATURE_INFRARED_ADDRESS, uint32(math.Round(h.TargetTemperature*10))))
		controller.AddAction(triggerAction)
		break
	case DURATION_JUDGE_TYPE:
		delayAction := action.NewDelayAction(h.Duration * 1000)
		controller.AddAction(delayAction)
	default:
		delayAction := action.NewDelayAction(1000)
		controller.AddAction(delayAction)
		logger.Log.Println("wrong temperature judge type")
	}
}

type DishOutInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDishOutInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(DISH_OUT),
	}
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
}

type ShakeInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewShakeInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(SHAKE),
	}
}

func (s ShakeInstruction) AddToController(controller *core.Controller) {
	shakeControlAction := action.NewShakeControlAction(data.SHAKE_CONTROL_WORD_ADDRESS,
		data.SHAKE_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.SHAKE_AMOUNT_ADDRESS, 5),
		data.NewAddressValue(data.SHAKE_UPWARD_SPEED_ADDRESS, 30000),
		data.NewAddressValue(data.SHAKE_DOWNWARD_SPEED_ADDRESS, 20000),
		data.NewAddressValue(data.SHAKE_DISTANCE_ADDRESS, 2000))
	controller.AddAction(shakeControlAction)
}

const (
	VENTING uint32 = iota + 1
	PURIFICATION
)

type LampblackPurifyInstruction struct {
	Instruction `mapstructure:",squash"`
	Mode        uint32 `json:"mode"`
}

func NewLampblackPurifyInstruction(mode uint32) *LampblackPurifyInstruction {
	return &LampblackPurifyInstruction{
		Instruction: NewInstruction(LAMPBLACK_PURIFY),
		Mode:        mode,
	}
}

func (l LampblackPurifyInstruction) AddToController(controller *core.Controller) {
	lampblackPurifyControlAction := action.NewLampblackPurifyControlAction(data.LAMPBLACK_PURIFY_CONTROL_WORD_ADDRESS,
		data.LAMPBLACK_PURIFY_STATUS_WORD_ADDRESS,
		data.NewAddressValue(data.LAMPBLACK_PURIFY_MODE_ADDRESS, l.Mode))
	controller.AddAction(lampblackPurifyControlAction)
}

type DoorUnlockInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDoorUnlockInstruction() *DoorUnlockInstruction {
	return &DoorUnlockInstruction{
		Instruction: NewInstruction(DOOR_UNLOCK),
	}
}

func (d DoorUnlockInstruction) AddToController(controller *core.Controller) {
	doorUnlockControlAction := action.NewDoorUnlockControlAction(data.DOOR_UNLOCK_CONTROL_WORD_ADDRESS,
		data.DOOR_UNLOCK_STATUS_WORD_ADDRESS)
	controller.AddAction(doorUnlockControlAction)
}

var InstructionTypeToStruct = map[InstructionType]Instructioner{
	INGREDIENT:       IngredientInstruction{},
	SEASONING:        SeasoningInstruction{},
	WATER:            WaterInstruction{},
	STIR_FRY:         StirFryInstruction{},
	HEAT:             HeatingInstruction{},
	DISH_OUT:         DishOutInstruction{},
	SHAKE:            ShakeInstruction{},
	LAMPBLACK_PURIFY: LampblackPurifyInstruction{},
	DOOR_UNLOCK:      DoorUnlockInstruction{},

	AXIS:   AxisInstruction{},
	ROTATE: RotateInstruction{},
	PUMP:   PumpInstruction{},
}
