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
}

func NewTriggerAction(triggerAddressValue *data.AddressValue) *TriggerAction {
	return &TriggerAction{
		BaseAction:          newBaseAction(TRIGGER),
		TriggerAddressValue: triggerAddressValue,
	}
}

func (t *TriggerAction) Execute(writer *operator.Writer, reader *operator.Reader) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if reader.Trig(t.TriggerAddressValue.Address, t.TriggerAddressValue.Value) {
				return
			}
		}
	}
}

func (t *TriggerAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[等待]触发%s>=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
}

func (t *TriggerAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]触发%s>=%d", t.TriggerAddressValue.Address, t.TriggerAddressValue.Value)
}
