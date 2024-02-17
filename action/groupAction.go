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
	//Actions []Actioner
	Actions chan Actioner
	wg      *sync.WaitGroup

	triggerActions chan *TriggerAction
}

func NewGroupAction() *GroupAction {
	return &GroupAction{
		BaseAction: newBaseAction(GROUP),
		//Actions:    []Actioner{},
		Actions: make(chan Actioner, 100),
		wg:      new(sync.WaitGroup),

		triggerActions: make(chan *TriggerAction, 100),
	}
}

func (g *GroupAction) AddAction(actioner Actioner) {
	//g.Actions = append(g.Actions, actioner)
	g.Actions <- actioner
}

func (g *GroupAction) CheckType() ActionType {
	return g.ActionType
}

func (g *GroupAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	if len(g.Actions) == 0 {
		return
	}
	g.wg.Add(len(g.Actions))

	for action := range g.Actions {
		go func(action Actioner) {
			defer g.wg.Done()
			logger.Log.Println(action.BeforeExecuteInfo())
			action.Execute(writer, reader, debugMode)
			if action.CheckType() == CONTROL {
				triggerAction := NewTriggerAction(data.NewAddressValue(action.GetStatusWordAddress(), 100), data.EQUAL_TO_TARGET, action.GetControlWordAddress())
				g.triggerActions <- triggerAction
				triggerAction.Execute(writer, reader, debugMode)
				logger.Log.Println(triggerAction.AfterExecuteInfo())
			}
		}(action)
		if len(g.Actions) == 0 {
			close(g.Actions)
		}
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

func (g *GroupAction) Shutdown() {
	logger.Log.Printf("[终止]同步执行%d个动作", len(g.Actions))
	for triggerAction := range g.triggerActions {
		go triggerAction.Shutdown()
		if len(g.triggerActions) == 0 {
			close(g.triggerActions)
		}
	}
}
