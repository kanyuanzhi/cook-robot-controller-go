package data

type BaseAction struct {
	ActionType string
}

type DelayAction struct {
	*BaseAction

	Delay uint32
}

type TriggerAction struct {
	*BaseAction
}

type Action struct {
	ControlWordAddress string
	StatusWordAddress  string
}

func NewAction(controlWordAddress string, statusWordAddress string) *Action {
	return &Action{
		ControlWordAddress: controlWordAddress,
		StatusWordAddress:  statusWordAddress,
	}
}

// StopAction 停止动作
type StopAction struct {
	*Action
}

func NewStopAction(controlWordAddress string, statusWordAddress string) *StopAction {
	return &StopAction{
		Action: NewAction(controlWordAddress, statusWordAddress),
	}
}

// AxleResetAction 轴重置动作
type AxleResetAction struct {
	*Action
}

func NewAxleResetAction(controlWordAddress string, statusWordAddress string) *AxleResetAction {
	return &AxleResetAction{
		Action: NewAction(controlWordAddress, statusWordAddress),
	}
}
