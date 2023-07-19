package action

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/operator"
	"fmt"
)

//type ControlActioner interface {
//	StopControlAction | AxisResetControlAction | AxisLocateControlAction | AxisRotateControlAction |
//		DishOutControlAction | PumpControlAction | LampblackPurifyControlAction | DoorLockControlAction |
//		TemperatureControlAction
//}

// StopAction 停止动作，用于R1轴停转
type StopAction struct {
	*BaseAction
	ControlWordAddress string // 控制字
	StatusWordAddress  string // 状态字
	AddressValueList   []*data.AddressValue
}

func NewStopAction(controlWordAddress string, statusWordAddress string) *StopAction {
	return &StopAction{
		BaseAction:         newBaseAction(STOP),
		ControlWordAddress: controlWordAddress,
		StatusWordAddress:  statusWordAddress,
		AddressValueList:   []*data.AddressValue{data.NewAddressValue(controlWordAddress, 0)},
	}
}

func (s *StopAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	if debugMode {
		return
	}
	successChan := make(chan bool)
	defer close(successChan)
	go writer.Send(successChan, s.AddressValueList)
	<-successChan
	return
}

func (s *StopAction) GetStatusWordAddress() string {
	return s.StatusWordAddress
}

func (s *StopAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]停止,发送(%s,%d),状态字地址%s判定为0",
		s.ControlWordAddress, 0, s.StatusWordAddress)
}

func (s *StopAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]停止,发送(%s,%d),状态字地址%s判定为0",
		s.ControlWordAddress, 0, s.StatusWordAddress)
}
