package action

import (
	"cook-robot-controller-go/operator"
)

type Actioner interface {
	CheckType() ActionType
	Execute(writer *operator.Writer, reader *operator.Reader)

	GetStatusWordAddress() string

	BeforeExecuteInfo() string
	AfterExecuteInfo() string
}

type ActionType uint8 // 动作类型

const (
	DELAY   ActionType = iota + 1 // 延时动作
	TRIGGER                       // 触发动作
	CONTROL                       // 控制动作
)

type BaseAction struct {
	ActionType ActionType
}

func newBaseAction(actionType ActionType) *BaseAction {
	return &BaseAction{
		ActionType: actionType,
	}
}

func (b *BaseAction) CheckType() ActionType {
	return b.ActionType
}

func (b *BaseAction) Execute(writer *operator.Writer, reader *operator.Reader) {
	return
}

func (b *BaseAction) GetStatusWordAddress() string {
	return ""
}

func (b *BaseAction) BeforeExecuteInfo() string {
	return ""
}

func (b *BaseAction) AfterExecuteInfo() string {
	return ""
}
