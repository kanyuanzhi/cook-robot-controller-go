package action

import (
	"cook-robot-controller-go/operator"
	"fmt"
	"time"
)

type DelayAction struct {
	*BaseAction
	Delay int64 // 延时
}

func NewDelayAction(delay int64) *DelayAction {
	return &DelayAction{
		BaseAction: newBaseAction(DELAY),
		Delay:      delay,
	}
}

func (d *DelayAction) Execute(writer *operator.Writer, reader *operator.Reader) {
	timer := time.NewTimer(time.Duration(d.Delay) * time.Millisecond)
	defer timer.Stop()
	//exitFlag := false
	//for {
	//	select {
	//	case <-timer.C:
	//		finishChan <- true
	//		exitFlag = true
	//		logger.Log.Printf("延时%dms完毕", d.Delay)
	//	default:
	//		//logger.Log.Println("延时执行中...")
	//		time.Sleep(10 * time.Millisecond)
	//	}
	//	if exitFlag {
	//		logger.Log.Println("跳出循环")
	//		break
	//	}
	//}
	<-timer.C
	return
}

func (d *DelayAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]延时执行%.3f秒", float32(d.Delay)/1000)
}

func (d *DelayAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]延时执行%.3f秒", float32(d.Delay)/1000)
}
