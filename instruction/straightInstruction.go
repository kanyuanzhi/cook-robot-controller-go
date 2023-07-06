package instruction

type function string

const (
	RESET  = function("reset")
	LOCATE = function("locate")
)

type AlexInstruction struct {
	*Instruction
	Function       function `json:"function" mapstructure:"function"`
	Alex           string   `json:"alex"  mapstructure:"alex"`
	TargetPosition uint32   `json:"target_position"  mapstructure:"target_position"`
	Speed          uint32   `json:"speed" mapstructure:"speed"`
}

func NewAlexInstruction(function function, alex string, targetPosition uint32, speed uint32) *AlexInstruction {
	return &AlexInstruction{
		Instruction:    NewInstruction(ALEX),
		Function:       function,
		Alex:           alex,
		TargetPosition: targetPosition,
		Speed:          speed,
	}
}

type RotateInstruction struct {
	*Instruction
	Function         function `json:"function"`
	Mode             uint8    `json:"mode"`
	Speed            uint32   `json:"speed"`
	RotationalAmount uint8    `json:"rotational_amount"`
}

func NewRotateInstruction(function function, mode uint8, speed uint32, rotationalAmount uint8) *RotateInstruction {
	return &RotateInstruction{
		Instruction:      NewInstruction(ROTATE),
		Function:         function,
		Mode:             mode,
		Speed:            speed,
		RotationalAmount: rotationalAmount,
	}
}

type PumpInstruction struct {
	*Instruction
	PumpNumber uint8  `json:"pumpNumber"`
	Duration   uint32 `json:"duration"`
}

func NewPumpInstruction(pumpNumber uint8, duration uint32) *PumpInstruction {
	return &PumpInstruction{
		Instruction: NewInstruction(PUMP),
		PumpNumber:  pumpNumber,
		Duration:    duration,
	}
}
