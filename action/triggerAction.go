package action

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/operator"
	"fmt"
	"time"
)

type TriggerAction struct {
	*BaseAction
	TriggerAddressValue *data.AddressValue // 触发地址-值
	triggerType         data.TriggerType
}

func NewTriggerAction(triggerAddressValue *data.AddressValue, triggerType data.TriggerType) *TriggerAction {
	return &TriggerAction{
		BaseAction:          newBaseAction(TRIGGER),
		TriggerAddressValue: triggerAddressValue,
		triggerType:         triggerType,
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

func (t *TriggerAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[等待]触发%s=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
}

func (t *TriggerAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]触发%s=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
}
