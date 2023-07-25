package core

import (
	"cook-robot-controller-go/action"
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/modbus"
	"cook-robot-controller-go/operator"
	"time"
)

type Controller struct {
	writer    *operator.Writer
	reader    *operator.Reader
	TcpServer *modbus.TCPServer

	waitingActionChan      chan action.Actioner
	executingActionChan    chan action.Actioner
	executingAction        action.Actioner // 当前在执行的动作
	executedActionFlagChan chan bool
	lastExecutedAction     action.Actioner // 前一个执行的动作

	pauseChan chan bool

	CurrentCommandName string // 当前命令名称 cook|wash|prepare|...
	CurrentDishUuid    string // 当前正在炒制菜品的uuid

	InstructionInfoChan    chan *data.InstructionInfo // 等待运行的指令队列，每个指令包含若干动作
	CurrentInstructionInfo *data.InstructionInfo      // 当前在运行的指令
	InstructionFlagChan    chan bool                  // 长度与waitingActionChan相同，每个flag对应一个action，当该action是本指令的第一个action时flag为true，否则为false

	IsPausing                       bool // 暂停中，只在炒制菜品过程中可以暂停
	IsRunning                       bool // 运行中，包括炒制菜品、备菜、清洗、出菜等multiple命令执行中true，执行完毕false
	IsCooking                       bool // 炒制菜品中
	IsPausingWithMovingFinished     bool // 中途加料暂停过程，移动y轴至加料位中false，完毕true
	IsPausingWithMovingBackFinished bool // 中途加料恢复过程，移动y轴至原位置中false，完毕true
	IsPausePermitted                bool // 是否允许暂停，翻炒延时过程中才允许暂停中途加料

	CookingTime         int64 // 炒制菜品已运行时长
	pauseCookTimingChan chan bool
	pauseCookTimingFlag bool
	stopCookTimingChan  chan bool

	lastYAxisLocateControlAction action.Actioner

	debugMode            bool
	MaxActionNumber      int
	MaxInstructionNumber int
}

func NewController(writer *operator.Writer, reader *operator.Reader, tcpServer *modbus.TCPServer, debugMode bool) *Controller {
	maxActionNumber := 100
	maxInstructionNumber := 50
	controller := &Controller{
		writer:                          writer,
		reader:                          reader,
		TcpServer:                       tcpServer,
		waitingActionChan:               make(chan action.Actioner, maxActionNumber),
		executingActionChan:             make(chan action.Actioner),
		executedActionFlagChan:          make(chan bool, maxActionNumber),
		pauseChan:                       make(chan bool),
		CurrentCommandName:              "",
		CurrentDishUuid:                 "",
		InstructionInfoChan:             make(chan *data.InstructionInfo, maxInstructionNumber),
		CurrentInstructionInfo:          &data.InstructionInfo{},
		InstructionFlagChan:             make(chan bool, maxActionNumber),
		IsPausing:                       false,
		IsRunning:                       false,
		IsCooking:                       false,
		IsPausingWithMovingFinished:     true,
		IsPausingWithMovingBackFinished: true,
		IsPausePermitted:                false,
		CookingTime:                     0,
		pauseCookTimingChan:             make(chan bool),
		pauseCookTimingFlag:             false,
		stopCookTimingChan:              make(chan bool),
		lastYAxisLocateControlAction:    nil,

		debugMode:            debugMode,
		MaxActionNumber:      maxActionNumber,
		MaxInstructionNumber: maxInstructionNumber,
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
	c.executedActionFlagChan <- true
	if a.CheckType() == action.CONTROL {
		triggerAction := action.NewTriggerAction(data.NewAddressValue(a.GetStatusWordAddress(), 100),
			data.EQUAL_TO_TARGET, a.GetControlWordAddress())
		c.waitingActionChan <- triggerAction
		c.executedActionFlagChan <- true
	}
	if a.CheckType() == action.STOP {
		triggerAction := action.NewTriggerAction(data.NewAddressValue(a.GetStatusWordAddress(), 0),
			data.EQUAL_TO_TARGET, a.GetControlWordAddress())
		c.waitingActionChan <- triggerAction
		c.executedActionFlagChan <- true
	}
}

func (c *Controller) AddInstructionInfo(insInfo *data.InstructionInfo) {
	insInfo.Index = len(c.InstructionInfoChan) + 1 // 指令序号从1开始
	c.InstructionInfoChan <- insInfo
	for i := 0; i < insInfo.ActionNumber; i++ {
		if i == 0 {
			c.InstructionFlagChan <- true
		} else {
			c.InstructionFlagChan <- false
		}
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
	logger.Log.Printf("[开始运行]%s", c.CurrentCommandName)
	go c.startTiming()
	c.IsRunning = true
	quitFlag := false
	isNextInstruction := false
	for {
		select {
		case executingAction := <-c.waitingActionChan:
			isNextInstruction = <-c.InstructionFlagChan
			if isNextInstruction {
				c.CurrentInstructionInfo = <-c.InstructionInfoChan
			}
			go func() {
				if c.CurrentCommandName == data.COMMAND_NAME_COOK && executingAction.CheckType() == action.DELAY {
					// 炒制菜品且当前待执行动作为延时时，
					if c.lastExecutedAction != nil &&
						c.lastExecutedAction.CheckType() == action.TRIGGER &&
						c.lastExecutedAction.GetControlWordAddress() == data.Y_LOCATE_CONTROL_WORD_ADDRESS {
						// 上一个动作为trigger且由y轴定位动作触发，允许中途加料
						c.IsPausePermitted = true
					}
				}

				if executingAction.CheckType() != action.TRIGGER {
					logger.Log.Println(executingAction.BeforeExecuteInfo())
				}

				executingAction.Execute(c.writer, c.reader, c.debugMode)

				if executingAction.CheckType() == action.TRIGGER {
					logger.Log.Println(executingAction.AfterExecuteInfo())
				}

				c.IsPausePermitted = false // 延时结束后不再允许中途加料

				c.lastExecutedAction = <-c.executingActionChan
				<-c.executedActionFlagChan
				if len(c.executedActionFlagChan) == 0 {
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
	logger.Log.Printf("[结束运行]%s", c.CurrentCommandName)
	c.Stop()
}

func (c *Controller) Pause() {
	c.pauseTiming()
	c.IsPausing = true
	//if c.executingAction.CheckType() == action.DELAY {
	//	c.executingAction.Pause()
	//} else if c.executingAction.CheckType() == action.CONTROL {
	//	// control型动作需要在执行trig后才能暂停
	//	triggerAction := <-c.waitingActionChan // control型动作后一定跟着一个trig型动作
	//	triggerAction.Execute(c.writer, c.reader, c.debugMode)
	//	<-c.executedActionChan // 总action数减1
	//	c.pauseChan <- true
	//} else {
	//	c.pauseChan <- true
	//}
	if c.executingAction.CheckType() == action.DELAY {
		go c.executingAction.Pause()
	} else {
		c.pauseChan <- true
	}
}

func (c *Controller) Resume() {
	c.resumeTiming()
	c.IsPausing = false
	if c.executingAction.CheckType() == action.DELAY {
		c.executingAction.Resume()
	} else {
		c.pauseChan <- true
	}
	// 如果之前有做过y轴定位动作，需要在恢复运行前将y轴定位到暂停前的位置
	//c.IsPausingWithMovingBackFinished = false
	//if c.lastYAxisLocateControlAction != nil {
	//	c.lastYAxisLocateControlAction.Execute(c.writer, c.reader, c.debugMode)
	//	triggerAction := action.NewTriggerAction(data.NewAddressValue(data.Y_LOCATE_STATUS_WORD_ADDRESS, 100), data.EQUAL_TO_TARGET)
	//	triggerAction.Execute(c.writer, c.reader, c.debugMode)
	//}
	//c.IsPausingWithMovingBackFinished = true
	//c.pauseChan <- true
}

func (c *Controller) Stop() {
	c.stopTiming()

	c.lastExecutedAction = nil

	c.waitingActionChan = make(chan action.Actioner, c.MaxActionNumber)
	c.CurrentCommandName = ""
	c.CurrentDishUuid = ""

	c.IsRunning = false
	c.IsCooking = false

	c.IsPausePermitted = false
	c.CookingTime = 0

	c.lastYAxisLocateControlAction = nil
}

func (c *Controller) startTiming() {
	if c.CurrentCommandName != data.COMMAND_NAME_COOK {
		return
	}
	ticker := time.NewTicker(250 * time.Millisecond)
	c.pauseCookTimingFlag = false
	defer ticker.Stop()
	quitTimingFlag := false
	for {
		select {
		case <-ticker.C:
			if !c.pauseCookTimingFlag {
				c.CookingTime += 250
			}
		case p := <-c.pauseCookTimingChan:
			c.pauseCookTimingFlag = p
		case <-c.stopCookTimingChan:
			quitTimingFlag = true
		}
		if quitTimingFlag {
			break
		}
	}
}

func (c *Controller) pauseTiming() {
	if c.CurrentCommandName != data.COMMAND_NAME_COOK {
		return
	}
	c.pauseCookTimingChan <- true
}

func (c *Controller) resumeTiming() {
	if c.CurrentCommandName != data.COMMAND_NAME_COOK {
		return
	}
	c.pauseCookTimingChan <- false
}

func (c *Controller) stopTiming() {
	if c.CurrentCommandName != data.COMMAND_NAME_COOK {
		return
	}
	c.stopCookTimingChan <- true
	c.CookingTime = 0
}
