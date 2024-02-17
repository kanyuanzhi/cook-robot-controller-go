package instruction

import (
	"cook-robot-controller-go/core"
)

type Instructioner interface {
	CheckType() InstructionType
	SetName(name string)
	AddToController(controller *core.Controller)
	ExecuteImmediately(controller *core.Controller)
}

type InstructionType string

const (
	INGREDIENT       = InstructionType("ingredient")
	SEASONING        = InstructionType("seasoning")
	WATER            = InstructionType("water")
	OIL              = InstructionType("oil")
	STIR_FRY         = InstructionType("stir_fry")
	HEAT             = InstructionType("heat")
	DISH_OUT         = InstructionType("dish_out")
	SHAKE            = InstructionType("shake")
	LAMPBLACK_PURIFY = InstructionType("lampblack_purify")
	DOOR_UNLOCK      = InstructionType("door_unlock")
	RESET_XYT        = InstructionType("reset_xyt")
	RESET_RT         = InstructionType("reset_rt")
	PREPARE          = InstructionType("prepare")
	DELAY            = InstructionType("delay")
	RESUME           = InstructionType("resume")
	PAUSE_TO_ADD     = InstructionType("pause_to_add")
	WASH             = InstructionType("wash")
	POUR             = InstructionType("pour")
	INIT             = InstructionType("init")
	FINISH           = InstructionType("finish")
	WITHDRAW         = InstructionType("withdrawer")
	SHUTDOWN         = InstructionType("shutdown")

	OPEN_WATER_SOURCE_VALVE  = InstructionType("open_water_source_valve")
	CLOSE_WATER_SOURCE_VALVE = InstructionType("close_water_source_valve")
	OPEN_WATER_PUMP_VALVE    = InstructionType("open_water_pump_valve")
	CLOSE_WATER_PUMP_VALVE   = InstructionType("close_water_pump_valve")
	OPEN_NOZZLE_VALVE        = InstructionType("open_nozzle_valve")
	CLOSE_NOZZLE_VALVE       = InstructionType("close_nozzle_valve")
	OPEN_PUMP_7_VALVE        = InstructionType("open_pump_7_valve")
	CLOSE_PUMP_7_VALVE       = InstructionType("close_pump_7_valve")

	AXIS   = InstructionType("axis")
	ROTATE = InstructionType("rotate")
	PUMP   = InstructionType("pump")
)

type Instruction struct {
	InstructionType InstructionType `json:"instructionType" mapstructure:"instructionType"`
	InstructionName string          `json:"instructionName" mapstructure:"instructionName"`
}

func (ins Instruction) CheckType() InstructionType {
	return ins.InstructionType
}

func (ins Instruction) SetName(name string) {
	ins.InstructionName = name
}

func (ins Instruction) ExecuteImmediately(controller *core.Controller) {
	return
}

func NewInstruction(instructionType InstructionType) Instruction {
	return Instruction{InstructionType: instructionType}
}

type IngredientInstruction struct {
	Instruction `mapstructure:",squash"`
	SlotNumber  string `json:"slotNumber" mapstructure:"slotNumber"`
}

func NewIngredientInstruction(slotNumber string) *IngredientInstruction {
	return &IngredientInstruction{
		Instruction: NewInstruction(INGREDIENT),
		SlotNumber:  slotNumber,
	}
}

type SeasoningInstruction struct {
	Instruction     `mapstructure:",squash"`
	PumpToWeightMap map[string]uint32 `json:"pumpToWeightMap" mapstructure:"pumpToWeightMap"` // 泵号:重量
	PumpToRatioMap  map[string]uint32 `json:"pumpToRatioMap" mapstructure:"pumpToRatioMap"`   // 泵号:重量g与时长ms比例
}

func NewSeasoningInstruction(name string, pumpToWeightMap map[string]uint32, pumpToRatioMap map[string]uint32) *SeasoningInstruction {
	return &SeasoningInstruction{
		Instruction: Instruction{
			InstructionType: SEASONING,
			InstructionName: name,
		},
		PumpToWeightMap: pumpToWeightMap,
		PumpToRatioMap:  pumpToRatioMap,
	}
}

type WaterInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pumpNumber" mapstructure:"pumpNumber"`
	Weight      uint32 `json:"weight"`
	Ratio       uint32 `json:"ratio"`
}

func NewWaterInstruction(pumpNumber uint32, weight uint32, ratio uint32) *WaterInstruction {
	return &WaterInstruction{
		Instruction: NewInstruction(WATER),
		PumpNumber:  pumpNumber,
		Weight:      weight,
		Ratio:       ratio,
	}
}

type OilInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pumpNumber" mapstructure:"pumpNumber"`
	Weight      uint32 `json:"weight"`
	Ratio       uint32 `json:"ratio"`
}

func NewOilInstruction(pumpNumber uint32, weight uint32, ratio uint32) *OilInstruction {
	return &OilInstruction{
		Instruction: NewInstruction(OIL),
		PumpNumber:  pumpNumber,
		Weight:      weight,
		Ratio:       ratio,
	}
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

type HeatInstruction struct {
	Instruction       `mapstructure:",squash"`
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"targetTemperature" mapstructure:"targetTemperature"`
	Duration          uint32  `json:"duration"`
	JudgeType         uint    `json:"judgeType" mapstructure:"judgeType"`
}

func NewHeatInstruction(temperature float64, targetTemperature float64, duration uint32, judgeType uint) *HeatInstruction {
	return &HeatInstruction{
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
	NO_JUDGE
)

type DishOutInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDishOutInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(DISH_OUT),
	}
}

type ShakeInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewShakeInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(SHAKE),
	}
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

type DoorUnlockInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDoorUnlockInstruction() *DoorUnlockInstruction {
	return &DoorUnlockInstruction{
		Instruction: NewInstruction(DOOR_UNLOCK),
	}
}

type ResetXYTInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResetXYTInstruction() *ResetXYTInstruction {
	return &ResetXYTInstruction{
		Instruction: NewInstruction(RESET_XYT),
	}
}

type ResetRTInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResetRTInstruction(name string) *ResetRTInstruction {
	return &ResetRTInstruction{
		Instruction: Instruction{
			InstructionType: RESET_RT,
			InstructionName: name,
		},
	}
}

// InitInstruction 炒菜开始前准备动作
type InitInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewInitInstruction(name string) *InitInstruction {
	return &InitInstruction{
		Instruction: Instruction{
			InstructionType: INIT,
			InstructionName: name,
		},
	}
}

// FinishInstruction 炒菜结束后停止动作
type FinishInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewFinishInstruction(name string) *FinishInstruction {
	return &FinishInstruction{
		Instruction: Instruction{
			InstructionType: FINISH,
			InstructionName: name,
		},
	}
}

type PrepareInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewPrepareInstruction() *PrepareInstruction {
	return &PrepareInstruction{
		Instruction: NewInstruction(PREPARE),
	}
}

type WashInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewWashInstruction() *WashInstruction {
	return &WashInstruction{
		Instruction: NewInstruction(WASH),
	}
}

type PourInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewPourInstruction() *PourInstruction {
	return &PourInstruction{
		Instruction: NewInstruction(POUR),
	}
}

type DelayInstruction struct {
	Instruction `mapstructure:",squash"`
	Duration    uint32 `json:"duration"`
}

func NewDelayInstruction(duration uint32) *DelayInstruction {
	return &DelayInstruction{
		Instruction: NewInstruction(DELAY),
		Duration:    duration,
	}
}

type ResumeInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResumeInstruction() *ResumeInstruction {
	return &ResumeInstruction{
		Instruction: NewInstruction(RESUME),
	}
}

type PauseToAddInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewPauseToAddInstruction() *PauseToAddInstruction {
	return &PauseToAddInstruction{
		Instruction: NewInstruction(PAUSE_TO_ADD),
	}
}

type WithdrawInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewWithdrawInstruction() *WithdrawInstruction {
	return &WithdrawInstruction{
		Instruction: NewInstruction(WITHDRAW),
	}
}

type ShutdownInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewShutdownInstruction() *ShutdownInstruction {
	return &ShutdownInstruction{
		Instruction: NewInstruction(SHUTDOWN),
	}
}

type OpenWaterSourceValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewOpenWaterSourceValveInstruction() *OpenWaterSourceValveInstruction {
	return &OpenWaterSourceValveInstruction{
		Instruction: NewInstruction(OPEN_WATER_SOURCE_VALVE),
	}
}

type CloseWaterSourceValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewCloseWaterSourceValveInstruction() *CloseWaterSourceValveInstruction {
	return &CloseWaterSourceValveInstruction{
		Instruction: NewInstruction(CLOSE_WATER_SOURCE_VALVE),
	}
}

type OpenWaterPumpValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewOpenWaterPumpValveInstruction() *OpenWaterPumpValveInstruction {
	return &OpenWaterPumpValveInstruction{
		Instruction: NewInstruction(OPEN_WATER_PUMP_VALVE),
	}
}

type CloseWaterPumpValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewCloseWaterPumpValveInstruction() *CloseWaterPumpValveInstruction {
	return &CloseWaterPumpValveInstruction{
		Instruction: NewInstruction(CLOSE_WATER_PUMP_VALVE),
	}
}

type OpenNozzleValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewOpenNozzleValveInstruction() *OpenNozzleValveInstruction {
	return &OpenNozzleValveInstruction{
		Instruction: NewInstruction(OPEN_NOZZLE_VALVE),
	}
}

type CloseNozzleValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewCloseNozzleValveInstruction() *CloseNozzleValveInstruction {
	return &CloseNozzleValveInstruction{
		Instruction: NewInstruction(CLOSE_NOZZLE_VALVE),
	}
}

type OpenPump7ValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewOpenPump7ValveInstruction() *OpenPump7ValveInstruction {
	return &OpenPump7ValveInstruction{
		Instruction: NewInstruction(OPEN_PUMP_7_VALVE),
	}
}

type ClosePump7ValveInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewClosePump7ValveInstruction() *ClosePump7ValveInstruction {
	return &ClosePump7ValveInstruction{
		Instruction: NewInstruction(CLOSE_PUMP_7_VALVE),
	}
}

var InstructionTypeToStruct = map[InstructionType]Instructioner{
	INGREDIENT:       IngredientInstruction{},
	SEASONING:        SeasoningInstruction{},
	WATER:            WaterInstruction{},
	OIL:              OilInstruction{},
	STIR_FRY:         StirFryInstruction{},
	HEAT:             HeatInstruction{},
	DISH_OUT:         DishOutInstruction{},
	SHAKE:            ShakeInstruction{},
	LAMPBLACK_PURIFY: LampblackPurifyInstruction{},
	DOOR_UNLOCK:      DoorUnlockInstruction{},
	INIT:             InitInstruction{},
	FINISH:           FinishInstruction{},
	RESET_XYT:        ResetXYTInstruction{},
	RESET_RT:         ResetRTInstruction{},
	PREPARE:          PrepareInstruction{},
	WASH:             WashInstruction{},
	POUR:             PourInstruction{},
	DELAY:            DelayInstruction{},
	RESUME:           ResumeInstruction{},
	PAUSE_TO_ADD:     PauseToAddInstruction{},
	WITHDRAW:         WithdrawInstruction{},
	SHUTDOWN:         ShutdownInstruction{},

	OPEN_WATER_SOURCE_VALVE:  OpenWaterSourceValveInstruction{},
	CLOSE_WATER_SOURCE_VALVE: CloseWaterSourceValveInstruction{},
	OPEN_WATER_PUMP_VALVE:    OpenWaterPumpValveInstruction{},
	CLOSE_WATER_PUMP_VALVE:   CloseWaterPumpValveInstruction{},
	OPEN_NOZZLE_VALVE:        OpenNozzleValveInstruction{},
	CLOSE_NOZZLE_VALVE:       CloseNozzleValveInstruction{},
	OPEN_PUMP_7_VALVE:        OpenPump7ValveInstruction{},
	CLOSE_PUMP_7_VALVE:       ClosePump7ValveInstruction{},

	AXIS:   AxisInstruction{},
	ROTATE: RotateInstruction{},
	PUMP:   PumpInstruction{},
}
