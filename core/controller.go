package core

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/operator"
)

type Controller struct {
	writer              *operator.Writer
	reader              *operator.Reader
	waitingActionChan   chan action.Actioner
	executingActionChan chan action.Actioner
	executedActionChan  chan bool
}

func NewController(writer *operator.Writer, reader *operator.Reader) *Controller {
	maxActionsNumber := 100
	controller := &Controller{
		writer:              writer,
		reader:              reader,
		waitingActionChan:   make(chan action.Actioner, maxActionsNumber),
		executingActionChan: make(chan action.Actioner),
		executedActionChan:  make(chan bool, maxActionsNumber),
	}
	return controller
}

func (c *Controller) Run() {

}

func (c *Controller) AddAction(a action.Actioner) {
	c.waitingActionChan <- a
	c.executedActionChan <- true
	if a.CheckType() == action.CONTROL {
		triggerAction := action.NewTriggerAction(data.NewAddressValue(a.GetStatusWordAddress(), 100))
		c.waitingActionChan <- triggerAction
		c.executedActionChan <- true
	}
}

func (c *Controller) Start() {
	logger.Log.Println("开始运行")

	for executingAction := range c.waitingActionChan {
		go func() {
			if executingAction.CheckType() != action.TRIGGER {
				logger.Log.Println(executingAction.BeforeExecuteInfo())
			}
			executingAction.Execute(c.writer, c.reader)
			if executingAction.CheckType() == action.TRIGGER {
				logger.Log.Println(executingAction.AfterExecuteInfo())
			}
			<-c.executingActionChan
			<-c.executedActionChan
			if len(c.executedActionChan) == 0 {
				// 所有action执行完毕，关闭waitingActionChan跳出for循环
				close(c.waitingActionChan)
			}
		}()
		c.executingActionChan <- executingAction
	}
	logger.Log.Println("结束运行")
	c.waitingActionChan = make(chan action.Actioner, 100)
}

func (c *Controller) Pause() {

}

func (c *Controller) Stop() {

}
