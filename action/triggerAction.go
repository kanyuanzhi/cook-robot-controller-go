package action

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/logger"
	"cook-robot-controller-go/operator"
	"fmt"
	"time"
)

type TriggerAction struct {
	*BaseAction
	TriggerAddressValue       *data.AddressValue // 触发地址-值
	triggerType               data.TriggerType
	TriggerControlWordAddress string // 使用触发动作的控制字地址
}

func NewTriggerAction(triggerAddressValue *data.AddressValue, triggerType data.TriggerType, triggerControlWordAddress string) *TriggerAction {
	return &TriggerAction{
		BaseAction:                newBaseAction(TRIGGER),
		TriggerAddressValue:       triggerAddressValue,
		triggerType:               triggerType,
		TriggerControlWordAddress: triggerControlWordAddress,
	}
}

func (t *TriggerAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	if debugMode {
		time.Sleep(1 * time.Second)
		return
	}
	time.Sleep(200 * time.Millisecond) // 延时200ms执行trig，确保状态字重置
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if reader.Trig(t.TriggerAddressValue.Address, t.TriggerAddressValue.Value, t.triggerType) {
				return
			}
		}
	}
}

func (t *TriggerAction) GetControlWordAddress() string {
	return t.TriggerControlWordAddress
}

func (t *TriggerAction) BeforeExecuteInfo() string {
	if t.triggerType == data.EQUAL_TO_TARGET {
		return fmt.Sprintf("[等待]触发%s=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
	} else if t.triggerType == data.LARGER_THAN_TARGET {
		return fmt.Sprintf("[等待]触发%s>=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
	} else {
		logger.Log.Println("触发条件错误")
		return ""
	}
}

func (t *TriggerAction) AfterExecuteInfo() string {
	if t.triggerType == data.EQUAL_TO_TARGET {
		return fmt.Sprintf("[结束]触发%s=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
	} else if t.triggerType == data.LARGER_THAN_TARGET {
		return fmt.Sprintf("[结束]触发%s>=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
	} else {
		logger.Log.Println("触发条件错误")
		return ""
	}
}
