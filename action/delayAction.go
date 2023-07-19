package action

import (
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/operator"
	"fmt"
	"time"
)

type DelayAction struct {
	*BaseAction
	Delay uint32 // 延时，毫秒

	pauseChan  chan bool
	resumeChan chan bool
}

func NewDelayAction(delay uint32) *DelayAction {
	return &DelayAction{
		BaseAction: newBaseAction(DELAY),
		Delay:      delay,
		pauseChan:  make(chan bool),
		resumeChan: make(chan bool),
	}
}

func (d *DelayAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	delayTime := d.Delay
	timer := time.NewTimer(time.Duration(delayTime) * time.Millisecond)
	beginTime := time.Now()
	defer timer.Stop()
	exitFlag := false
	var remainingTime uint32
	for {
		select {
		case <-timer.C:
			//finishChan <- true
			exitFlag = true
		//logger.Log.Printf("延时%dms完毕", d.Delay)
		case <-d.pauseChan:
			if !timer.Stop() {
				<-timer.C
			}
			passedTime := time.Until(beginTime)
			remainingTime = delayTime + uint32(passedTime.Milliseconds())
			logger.Log.Printf("[计时暂停]剩余时间%.3fs", float32(remainingTime)/1000)
		case <-d.resumeChan:
			logger.Log.Println("[计时恢复]")
			timer.Reset(time.Duration(remainingTime) * time.Millisecond)
			beginTime = time.Now()
			delayTime = remainingTime
		}
		if exitFlag {
			break
		}
	}
	//<-timer.C
	return
}

func (d *DelayAction) Pause() {
	d.pauseChan <- true
}

func (d *DelayAction) Resume() {
	d.resumeChan <- true
}

func (d *DelayAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]延时执行%.3f秒", float32(d.Delay)/1000)
}

func (d *DelayAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]延时执行%.3f秒", float32(d.Delay)/1000)
}
