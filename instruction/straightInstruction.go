package instruction

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"strings"
)

type function string

const (
	RESET  = function("reset")
	LOCATE = function("locate")
	START  = function("start")
	STOP   = function("stop")
)

type AxisInstruction struct {
	Instruction    `mapstructure:",squash"`
	Function       function `json:"function"`
	Axis           string   `json:"AXIS"`
	TargetPosition uint32   `json:"target_position"  mapstructure:"target_position"`
	Speed          uint32   `json:"speed"`
}

func NewAxisInstruction(function function, axis string, targetPosition uint32, speed uint32) *AxisInstruction {
	return &AxisInstruction{
		Instruction:    NewInstruction(AXIS),
		Function:       function,
		Axis:           axis,
		TargetPosition: targetPosition,
		Speed:          speed,
	}
}

func (a AxisInstruction) AddToController(controller *core.Controller) {
	if a.Function != RESET && a.Function != LOCATE {
		logger.Log.Println("wrong axis function")
		return
	}
	a.Axis = strings.ToUpper(a.Axis)
	if a.Axis != "X" && a.Axis != "Y" && a.Axis != "Z" && a.Axis != "R1" && a.Axis != "R2" {
		logger.Log.Println("wrong axis")
		return
	}
	var axisControlAction action.Actioner
	if a.Function == RESET { // 复位
		axisControlAction = action.NewAxisResetControlAction(data.AxisToResetControlWordAddress[a.Axis],
			data.AxisToResetStatusWordAddress[a.Axis])
	} else { // 定位
		axisControlAction = action.NewAxisLocateControlAction(data.AxisToLocateControlWordAddress[a.Axis],
			data.AxisToLocateStatusWordAddress[a.Axis],
			data.NewAddressValue(data.AxisToLocatePositionAddress[a.Axis], a.TargetPosition),
			data.NewAddressValue(data.AxisToLocateSpeedAddress[a.Axis], a.Speed))
	}
	controller.AddAction(axisControlAction)
}

type RotateInstruction struct {
	Instruction      `mapstructure:",squash"`
	Function         function `json:"function"`
	Mode             uint32   `json:"mode"`
	Speed            uint32   `json:"speed"`
	RotationalAmount uint32   `json:"rotational_amount" mapstructure:"rotational_amount"`
}

func NewRotateInstruction(function function, mode uint32, speed uint32, rotationalAmount uint32) *RotateInstruction {
	return &RotateInstruction{
		Instruction:      NewInstruction(ROTATE),
		Function:         function,
		Mode:             mode,
		Speed:            speed,
		RotationalAmount: rotationalAmount,
	}
}

func (r RotateInstruction) AddToController(controller *core.Controller) {
	if r.Function != STOP && r.Function != START {
		logger.Log.Println("wrong ratate function")
		return
	}
	var rotateControlAction action.Actioner
	if r.Function == STOP { // 停转
		rotateControlAction = action.NewStopAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
			data.R1_ROTATE_STATUS_WORD_ADDRESS)
	} else { // 定位
		rotateControlAction = action.NewAxisRotateControlAction(data.R1_ROTATE_CONTROL_WORD_ADDRESS,
			data.R1_ROTATE_STATUS_WORD_ADDRESS,
			data.NewAddressValue(data.R1_ROTATE_MODE_ADDRESS, r.Mode),
			data.NewAddressValue(data.R1_ROTATE_SPEED_ADDRESS, r.Mode),
			data.NewAddressValue(data.R1_ROTATE_AMOUNT_ADDRESS, r.Mode))
	}
	controller.AddAction(rotateControlAction)
}

type PumpInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint8  `json:"pump_number" mapstructure:"pump_number"`
	Duration    uint32 `json:"duration"`
}

func NewPumpInstruction(pumpNumber uint8, duration uint32) *PumpInstruction {
	return &PumpInstruction{
		Instruction: NewInstruction(PUMP),
		PumpNumber:  pumpNumber,
		Duration:    duration,
	}
}
