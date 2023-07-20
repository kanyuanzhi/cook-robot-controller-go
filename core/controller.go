package core

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/operator"
	"time"
)

type Controller struct {
	writer              *operator.Writer
	reader              *operator.Reader
	waitingActionChan   chan action.Actioner
	executingActionChan chan action.Actioner
	executingAction     action.Actioner
	executedActionChan  chan bool

	pauseChan chan bool

	CurrentCommandName  string
	PreviousCommandName string

	IsPausing                   bool // 暂停中，只在炒制菜品过程中可以暂停
	IsRunning                   bool // 运行中，包括炒制菜品、备菜、清洗、出菜等multiple命令执行中true，执行完毕false
	IsCooking                   bool // 炒制菜品中
	IsPausingWithMovingFinished bool // 中途加料暂停过程，移动y轴至加料位中false，完毕true

	lastYAxisLocateControlAction chan action.Actioner

	debugMode bool
}

func NewController(writer *operator.Writer, reader *operator.Reader, debugMode bool) *Controller {
	maxActionsNumber := 100
	controller := &Controller{
		writer:                       writer,
		reader:                       reader,
		waitingActionChan:            make(chan action.Actioner, maxActionsNumber),
		executingActionChan:          make(chan action.Actioner),
		executedActionChan:           make(chan bool, maxActionsNumber),
		pauseChan:                    make(chan bool),
		CurrentCommandName:           "",
		IsPausing:                    false,
		IsRunning:                    false,
		IsCooking:                    false,
		IsPausingWithMovingFinished:  false,
		lastYAxisLocateControlAction: make(chan action.Actioner),
		debugMode:                    debugMode,
	}
	return controller
}

func (c *Controller) Run() {
	if c.debugMode {
		logger.Log.Println("控制器以测试模式启动，判定动作延时1s完成，延时动作正常完成，其他动作立即完成")
		return
	}
	for {
		//logger.Log.Println(c.CurrentCommandName)
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) AddAction(a action.Actioner) {
	c.waitingActionChan <- a
	c.executedActionChan <- true
	if a.CheckType() == action.CONTROL {
		triggerAction := action.NewTriggerAction(data.NewAddressValue(a.GetStatusWordAddress(), 100), data.EQUAL_TO_TARGET) // 100
		c.waitingActionChan <- triggerAction
		c.executedActionChan <- true
	}
	if a.CheckType() == action.STOP {
		triggerAction := action.NewTriggerAction(data.NewAddressValue(a.GetStatusWordAddress(), 0), data.EQUAL_TO_TARGET)
		c.waitingActionChan <- triggerAction
		c.executedActionChan <- true
	}
}

func (c *Controller) ExecuteImmediately(a action.Actioner) {
	if a.CheckType() != action.TRIGGER {
		logger.Log.Println("[立即执行]", a.BeforeExecuteInfo())
	}
	a.Execute(c.writer, c.reader, c.debugMode)
	if a.CheckType() == action.TRIGGER {
		logger.Log.Println("[立即执行]", a.AfterExecuteInfo())
	}
}

func (c *Controller) Start() {
	logger.Log.Printf("[%s开始运行]", c.CurrentCommandName)
	c.IsRunning = true
	quitFlag := false
	for {
		select {
		case executingAction := <-c.waitingActionChan:
			go func() {
				//if executingAction.GetStatusWordAddress() == data.Y_LOCATE_STATUS_WORD_ADDRESS {
				//	if len(c.lastYAxisLocateControlAction) == 1 {
				//		<-c.lastYAxisLocateControlAction
				//	}
				//	c.lastYAxisLocateControlAction <- executingAction
				//}
				if executingAction.CheckType() != action.TRIGGER {
					logger.Log.Println(executingAction.BeforeExecuteInfo())
				}
				executingAction.Execute(c.writer, c.reader, c.debugMode)
				if executingAction.CheckType() == action.TRIGGER {
					logger.Log.Println(executingAction.AfterExecuteInfo())
				}
				<-c.executingActionChan
				<-c.executedActionChan
				if len(c.executedActionChan) == 0 {
					// 所有action执行完毕，关闭waitingActionChan跳出for循环
					close(c.waitingActionChan)
					quitFlag = true
				}
			}()
			c.executingAction = executingAction
			c.executingActionChan <- executingAction
		case <-c.pauseChan:
			<-c.pauseChan
		}

		if quitFlag {
			break
		}
	}
	logger.Log.Printf("[%s结束运行]", c.CurrentCommandName)
	c.waitingActionChan = make(chan action.Actioner, 100)
	c.CurrentCommandName = ""
	c.PreviousCommandName = ""

	c.IsRunning = false
	c.IsCooking = false
}

func (c *Controller) Pause() {
	c.IsPausing = true
	if c.executingAction.CheckType() == action.DELAY {
		c.executingAction.Pause()
	} else {
		c.pauseChan <- true
	}
	logger.Log.Println("暂停运行......")
	//c.PreviousCommandName = c.CurrentCommandName
	//c.CurrentCommandName = "pause"
}

func (c *Controller) Resume() {
	c.IsPausing = false
	if c.executingAction.CheckType() == action.DELAY {
		c.executingAction.Resume()
	} else {
		c.pauseChan <- true
	}
	//if len(c.lastYAxisLocateControlAction) == 1 {
	//	yAxisLocateControlAction := <-c.lastYAxisLocateControlAction
	//	yAxisLocateControlAction.Execute(c.writer, c.reader, c.debugMode)
	//}
	logger.Log.Println("恢复运行......")
	//c.CurrentCommandName = c.PreviousCommandName
}

func (c *Controller) Stop() {

}
