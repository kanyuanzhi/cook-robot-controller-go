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
	actions []Actioner
	wg      *sync.WaitGroup
}

func NewGroupAction() *GroupAction {
	return &GroupAction{
		BaseAction: newBaseAction(GROUP),
		actions:    []Actioner{},
		wg:         new(sync.WaitGroup),
	}
}

func (g *GroupAction) AddAction(actioner Actioner) {
	g.actions = append(g.actions, actioner)
}

func (g *GroupAction) CheckType() ActionType {
	return g.ActionType
}

func (g *GroupAction) Execute(writer *operator.Writer, reader *operator.Reader) {
	g.wg.Add(len(g.actions))

	for _, action := range g.actions {
		go func(action Actioner) {
			defer g.wg.Done()
			logger.Log.Println(action.BeforeExecuteInfo())
			action.Execute(writer, reader)
			if action.CheckType() == CONTROL {
				triggerAction := NewTriggerAction(data.NewAddressValue(action.GetStatusWordAddress(), 100))
				triggerAction.Execute(writer, reader)
				logger.Log.Println(triggerAction.AfterExecuteInfo())
			}
		}(action)
	}
	g.wg.Wait()

	return
}

func (g *GroupAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]同步执行%d个动作", len(g.actions))
}

func (g *GroupAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]同步执行%d个动作", len(g.actions))
}
