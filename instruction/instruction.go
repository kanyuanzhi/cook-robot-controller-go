package instruction

type Instructioner interface {
	CheckType() InstructionType
}

type InstructionType string

const (
	INGREDIENT      = InstructionType("ingredient")
	SEASONING       = InstructionType("seasoning")
	WATER           = InstructionType("water")
	STIR_FRY        = InstructionType("stir_fry")
	HEATING         = InstructionType("heating")
	DISH_OUT        = InstructionType("dish_out")
	LampblackPurify = InstructionType("lampblack_purify")
	DoorUnlock      = InstructionType("door_unlock")

	ALEX   = InstructionType("alex")
	ROTATE = InstructionType("rotate")
	PUMP   = InstructionType("pump")
)

type Instruction struct {
	InstructionType InstructionType `json:"instruction_type" mapstructure:"function"`
}

func (i *Instruction) CheckType() InstructionType {
	return i.InstructionType
}

func NewInstruction(instructionType InstructionType) *Instruction {
	return &Instruction{InstructionType: instructionType}
}

type IngredientInstruction struct {
	*Instruction
	SlotNumber uint8 `json:"slot_number"`
}

func NewIngredientInstruction(slotNumber uint8) *IngredientInstruction {
	return &IngredientInstruction{
		Instruction: NewInstruction(INGREDIENT),
		SlotNumber:  slotNumber,
	}
}

type SeasoningInstruction struct {
	*Instruction
	PumpNumber uint8  `json:"pump_number"`
	Weight     uint32 `json:"weight"`
}

func NewSeasoningInstruction(pumpNumber uint8, weight uint32) *SeasoningInstruction {
	return &SeasoningInstruction{
		Instruction: NewInstruction(SEASONING),
		PumpNumber:  pumpNumber,
		Weight:      weight,
	}
}

type WaterInstruction struct {
	*Instruction
	PumpNumber uint8  `json:"pump_number"`
	Weight     uint32 `json:"weight"`
}

func NewWaterInstruction(pumpNumber uint8, weight uint32) *WaterInstruction {
	return &WaterInstruction{
		Instruction: NewInstruction(WATER),
		PumpNumber:  pumpNumber,
		Weight:      weight,
	}
}

type StirFryInstruction struct {
	*Instruction
	Gear     uint8  `json:"gear"`
	Duration uint32 `json:"duration"`
}

func NewStirFryInstruction(gear uint8, duration uint32) *StirFryInstruction {
	return &StirFryInstruction{
		Instruction: NewInstruction(STIR_FRY),
		Gear:        gear,
		Duration:    duration,
	}
}

type HeatingInstruction struct {
	*Instruction
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"target_temperature"`
	Duration          uint32  `json:"duration"`
}

func NewHeatingInstruction(temperature float64, targetTemperature float64, duration uint32) *HeatingInstruction {
	return &HeatingInstruction{
		Instruction:       NewInstruction(HEATING),
		Temperature:       temperature,
		TargetTemperature: targetTemperature,
		Duration:          duration,
	}
}

type DishOutInstruction struct {
	*Instruction
}

func NewDishOutInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(DISH_OUT),
	}
}

type LampblackPurifyInstruction struct {
	*Instruction
}

func NewLampblackPurifyInstruction() *LampblackPurifyInstruction {
	return &LampblackPurifyInstruction{
		Instruction: NewInstruction(LampblackPurify),
	}
}

type DoorUnlockInstruction struct {
	*Instruction
}

func NewDoorUnlockInstruction() *DoorUnlockInstruction {
	return &DoorUnlockInstruction{
		Instruction: NewInstruction(DoorUnlock),
	}
}
