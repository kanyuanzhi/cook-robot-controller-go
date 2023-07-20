package action

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/operator"
	"fmt"
	"sync"
)

type GroupAction struct {
	*BaseAction
	Actions []Actioner
	wg      *sync.WaitGroup
}

func NewGroupAction() *GroupAction {
	return &GroupAction{
		BaseAction: newBaseAction(GROUP),
		Actions:    []Actioner{},
		wg:         new(sync.WaitGroup),
	}
}

func (g *GroupAction) AddAction(actioner Actioner) {
	g.Actions = append(g.Actions, actioner)
}

func (g *GroupAction) CheckType() ActionType {
	return g.ActionType
}

func (g *GroupAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	if len(g.Actions) == 0 {
		return
	}
	g.wg.Add(len(g.Actions))

	for _, action := range g.Actions {
		go func(action Actioner) {
			defer g.wg.Done()
			logger.Log.Println(action.BeforeExecuteInfo())
			action.Execute(writer, reader, debugMode)
			if action.CheckType() == CONTROL {
				triggerAction := NewTriggerAction(data.NewAddressValue(action.GetStatusWordAddress(), 100), data.EQUAL_TO_TARGET)
				triggerAction.Execute(writer, reader, debugMode)
				logger.Log.Println(triggerAction.AfterExecuteInfo())
			}
		}(action)
	}
	g.wg.Wait()

	return
}

func (g *GroupAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]同步执行%d个动作", len(g.Actions))
}

func (g *GroupAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]同步执行%d个动作", len(g.Actions))
}
