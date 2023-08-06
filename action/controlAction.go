package action

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/operator"
	"fmt"
)

// ControlAction 控制动作基类
type ControlAction struct {
	*BaseAction
	ControlWordAddress string // 控制字
	StatusWordAddress  string // 状态字
	AddressValueList   []*data.AddressValue
}

func NewControlAction(controlWordAddress string, statusWordAddress string) *ControlAction {
	return &ControlAction{
		BaseAction:         newBaseAction(CONTROL),
		ControlWordAddress: controlWordAddress,
		StatusWordAddress:  statusWordAddress,
		AddressValueList:   []*data.AddressValue{data.NewAddressValue(controlWordAddress, 1)},
	}
}

func (c *ControlAction) Execute(writer *operator.Writer, reader *operator.Reader, debugMode bool) {
	if debugMode {
		return
	}
	successChan := make(chan bool)
	defer close(successChan)
	go writer.Send(successChan, c.AddressValueList)
	<-successChan
	return
}

func (c *ControlAction) GetStatusWordAddress() string {
	return c.StatusWordAddress
}

func (c *ControlAction) GetControlWordAddress() string {
	return c.ControlWordAddress
}

// AxisResetControlAction 轴重置动作
type AxisResetControlAction struct {
	*ControlAction
}

func NewAxisResetControlAction(controlWordAddress string, statusWordAddress string) *AxisResetControlAction {
	axleResetControlAction := &AxisResetControlAction{
		ControlAction: NewControlAction(controlWordAddress, statusWordAddress),
	}
	return axleResetControlAction
}

func (a *AxisResetControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]%s轴复位,发送(%s,%d),状态字地址%s",
		data.AxisControlWordAddressToAxis[a.ControlWordAddress], a.ControlWordAddress, 1, a.StatusWordAddress)
}

func (a *AxisResetControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]%s轴复位,发送(%s,%d),状态字地址%s",
		data.AxisControlWordAddressToAxis[a.ControlWordAddress], a.ControlWordAddress, 1, a.StatusWordAddress)
}

// AxisLocateControlAction 轴定位动作
type AxisLocateControlAction struct {
	*ControlAction
	TargetPosition *data.AddressValue // 目标位置
	speed          *data.AddressValue // 移动速度
}

func NewAxisLocateControlAction(controlWordAddress string, statusWordAddress string,
	targetPosition *data.AddressValue, speed *data.AddressValue) *AxisLocateControlAction {
	axleLocateControlAction := &AxisLocateControlAction{
		ControlAction:  NewControlAction(controlWordAddress, statusWordAddress),
		TargetPosition: targetPosition,
		speed:          speed,
	}
	axleLocateControlAction.ControlAction.AddressValueList = append(axleLocateControlAction.ControlAction.AddressValueList,
		targetPosition, speed)
	return axleLocateControlAction
}

func (a *AxisLocateControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]%s轴定位%d,速度%d,发送(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		data.AxisControlWordAddressToAxis[a.ControlWordAddress], a.TargetPosition.Value, a.speed.Value,
		a.ControlWordAddress, 1,
		a.TargetPosition.Address, a.TargetPosition.Value,
		a.speed.Address, a.speed.Value,
		a.StatusWordAddress)
}

func (a *AxisLocateControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]%s轴定位%d,速度%d,发送(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		data.AxisControlWordAddressToAxis[a.ControlWordAddress], a.TargetPosition.Value, a.speed.Value,
		a.ControlWordAddress, 1,
		a.TargetPosition.Address, a.TargetPosition.Value,
		a.speed.Address, a.speed.Value,
		a.StatusWordAddress)
}

// AxisRotateControlAction 旋转动作
type AxisRotateControlAction struct {
	*ControlAction
	mode             *data.AddressValue // 旋转模式
	speed            *data.AddressValue // 旋转速度
	rotationalAmount *data.AddressValue // 正反转圈数
}

func NewAxisRotateControlAction(controlWordAddress string, statusWordAddress string,
	mode *data.AddressValue, speed *data.AddressValue, rotationalAmount *data.AddressValue) *AxisRotateControlAction {
	axleRotateControlAction := &AxisRotateControlAction{
		ControlAction:    NewControlAction(controlWordAddress, statusWordAddress),
		mode:             mode,
		speed:            speed,
		rotationalAmount: rotationalAmount,
	}
	axleRotateControlAction.ControlAction.AddressValueList = append(axleRotateControlAction.ControlAction.AddressValueList,
		mode, speed, rotationalAmount)
	return axleRotateControlAction
}

func (a *AxisRotateControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]R1轴%s,速度%d,正反转圈数%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		data.RotateModeToString[a.mode.Value], a.speed.Value, a.rotationalAmount.Value,
		a.ControlWordAddress, 1,
		a.mode.Address, a.mode.Value,
		a.speed.Address, a.speed.Value,
		a.rotationalAmount.Address, a.rotationalAmount.Value,
		a.StatusWordAddress)
}

func (a *AxisRotateControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]R1轴%s,速度%d,正反转圈数%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		data.RotateModeToString[a.mode.Value], a.speed.Value, a.rotationalAmount.Value,
		a.ControlWordAddress, 1,
		a.mode.Address, a.mode.Value,
		a.speed.Address, a.speed.Value,
		a.rotationalAmount.Address, a.rotationalAmount.Value,
		a.StatusWordAddress)
}

// DishOutControlAction 出菜动作
type DishOutControlAction struct {
	*ControlAction
	amount           *data.AddressValue // 抖动次数
	upwardSpeed      *data.AddressValue // 上行速度
	downwardSpeed    *data.AddressValue // 下行速度
	upwardPosition   *data.AddressValue // 上行位置
	downwardPosition *data.AddressValue // 下行位置
}

func NewDishOutControlAction(controlWordAddress string, statusWordAddress string,
	amount *data.AddressValue, upwardSpeed *data.AddressValue, downwardSpeed *data.AddressValue,
	upwardPosition *data.AddressValue, downwardPosition *data.AddressValue) *DishOutControlAction {
	dishOutControlAction := &DishOutControlAction{
		ControlAction:    NewControlAction(controlWordAddress, statusWordAddress),
		amount:           amount,
		upwardSpeed:      upwardSpeed,
		downwardSpeed:    downwardSpeed,
		upwardPosition:   upwardPosition,
		downwardPosition: downwardPosition,
	}
	dishOutControlAction.ControlAction.AddressValueList = append(dishOutControlAction.ControlAction.AddressValueList,
		amount, upwardSpeed, downwardSpeed, upwardPosition, downwardPosition)
	return dishOutControlAction
}

func (d *DishOutControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]一键出菜,抖动次数%d,上行速度%d,下行速度%d,出菜高位%d,出菜低位%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		d.amount.Value, d.upwardSpeed.Value, d.downwardSpeed.Value, d.upwardPosition.Value, d.downwardPosition.Value,
		d.ControlWordAddress, 1,
		d.amount.Address, d.amount.Value,
		d.upwardSpeed.Address, d.upwardSpeed.Value,
		d.downwardSpeed.Address, d.downwardSpeed.Value,
		d.upwardPosition.Address, d.upwardPosition.Value,
		d.downwardPosition.Address, d.downwardPosition.Value,
		d.StatusWordAddress)
}

func (d *DishOutControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]一键出菜,抖动次数%d,上行速度%d,下行速度%d,出菜高位%d,出菜低位%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		d.amount.Value, d.upwardSpeed.Value, d.downwardSpeed.Value, d.upwardPosition.Value, d.downwardPosition.Value,
		d.ControlWordAddress, 1,
		d.amount.Address, d.amount.Value,
		d.upwardSpeed.Address, d.upwardSpeed.Value,
		d.downwardSpeed.Address, d.downwardSpeed.Value,
		d.upwardPosition.Address, d.upwardPosition.Value,
		d.downwardPosition.Address, d.downwardPosition.Value,
		d.StatusWordAddress)
}

// ShakeControlAction 出菜动作
type ShakeControlAction struct {
	*ControlAction
	amount        *data.AddressValue // 抖动次数
	upwardSpeed   *data.AddressValue // 上行速度
	downwardSpeed *data.AddressValue // 下行速度
	distance      *data.AddressValue // 抖菜距离
}

func NewShakeControlAction(controlWordAddress string, statusWordAddress string,
	amount *data.AddressValue, upwardSpeed *data.AddressValue, downwardSpeed *data.AddressValue,
	distance *data.AddressValue) *ShakeControlAction {
	shakeControlAction := &ShakeControlAction{
		ControlAction: NewControlAction(controlWordAddress, statusWordAddress),
		amount:        amount,
		upwardSpeed:   upwardSpeed,
		downwardSpeed: downwardSpeed,
		distance:      distance,
	}
	shakeControlAction.ControlAction.AddressValueList = append(shakeControlAction.ControlAction.AddressValueList,
		amount, upwardSpeed, downwardSpeed, distance)
	return shakeControlAction
}

func (s *ShakeControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]菜盒抖菜,抖动次数%d,上行速度%d,下行速度%d,抖菜距离%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		s.amount.Value, s.upwardSpeed.Value, s.downwardSpeed.Value, s.distance.Value,
		s.ControlWordAddress, 1,
		s.amount.Address, s.amount.Value,
		s.upwardSpeed.Address, s.upwardSpeed.Value,
		s.downwardSpeed.Address, s.downwardSpeed.Value,
		s.distance.Address, s.distance.Value,
		s.StatusWordAddress)
}

func (s *ShakeControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]菜盒抖菜,抖动次数%d,上行速度%d,下行速度%d,抖菜距离%d,发送(%s,%d),(%s,%d),(%s,%d),(%s,%d),(%s,%d),状态字地址%s",
		s.amount.Value, s.upwardSpeed.Value, s.downwardSpeed.Value, s.distance.Value,
		s.ControlWordAddress, 1,
		s.amount.Address, s.amount.Value,
		s.upwardSpeed.Address, s.upwardSpeed.Value,
		s.downwardSpeed.Address, s.downwardSpeed.Value,
		s.distance.Address, s.distance.Value,
		s.StatusWordAddress)
}

// PumpControlAction 泵动作，液体泵1-6号、抽水泵7-8号、固体泵9-10号
type PumpControlAction struct {
	*ControlAction
	duration *data.AddressValue // 泵开启时长 ms
}

func NewPumpControlAction(controlWordAddress string, statusWordAddress string,
	duration *data.AddressValue) *PumpControlAction {
	pumpControlAction := &PumpControlAction{
		ControlAction: NewControlAction(controlWordAddress, statusWordAddress),
		duration:      duration,
	}
	pumpControlAction.ControlAction.AddressValueList = append(pumpControlAction.ControlAction.AddressValueList,
		duration)
	return pumpControlAction
}

func (p *PumpControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]%d号泵(%s)打开%d毫秒,发送(%s,%d),(%s,%d),状态字地址%s",
		data.PumpControlWordAddressToPumpNumber[p.ControlWordAddress],
		data.PumpNumberToPumpType[data.PumpControlWordAddressToPumpNumber[p.ControlWordAddress]], p.duration.Value,
		p.ControlWordAddress, 1,
		p.duration.Address, p.duration.Value,
		p.StatusWordAddress)
}

func (p *PumpControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]%d号泵(%s)打开%d毫秒,发送(%s,%d),(%s,%d),状态字地址%s",
		data.PumpControlWordAddressToPumpNumber[p.ControlWordAddress],
		data.PumpNumberToPumpType[data.PumpControlWordAddressToPumpNumber[p.ControlWordAddress]], p.duration.Value,
		p.ControlWordAddress, 1,
		p.duration.Address, p.duration.Value,
		p.StatusWordAddress)
}

// LampblackPurifyControlAction 油烟净化动作
type LampblackPurifyControlAction struct {
	*ControlAction
	mode *data.AddressValue // 模式,1:排气,2:排气+净化
}

func NewLampblackPurifyControlAction(controlWordAddress string, statusWordAddress string,
	mode *data.AddressValue) *LampblackPurifyControlAction {
	lampblackPurifyControlAction := &LampblackPurifyControlAction{
		ControlAction: NewControlAction(controlWordAddress, statusWordAddress),
		mode:          mode,
	}
	lampblackPurifyControlAction.ControlAction.AddressValueList = append(lampblackPurifyControlAction.ControlAction.AddressValueList,
		mode)
	return lampblackPurifyControlAction
}

func (l *LampblackPurifyControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]油烟净化,%s模式,发送(%s,%d),(%s,%d),状态字地址%s",
		data.LampblackPurifyModeToString[l.mode.Value],
		l.ControlWordAddress, 1,
		l.mode.Address, l.mode.Value,
		l.StatusWordAddress)
}

func (l *LampblackPurifyControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]油烟净化,%s模式,发送(%s,%d),(%s,%d),状态字地址%s",
		data.LampblackPurifyModeToString[l.mode.Value],
		l.ControlWordAddress, 1,
		l.mode.Address, l.mode.Value,
		l.StatusWordAddress)
}

// DoorUnlockControlAction 电磁门锁动作
type DoorUnlockControlAction struct {
	*ControlAction
}

func NewDoorUnlockControlAction(controlWordAddress string, statusWordAddress string) *DoorUnlockControlAction {
	doorLockControlAction := &DoorUnlockControlAction{
		ControlAction: NewControlAction(controlWordAddress, statusWordAddress),
	}
	return doorLockControlAction
}

func (d *DoorUnlockControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]电磁门解锁,发送(%s,%d),状态字地址%s",
		d.ControlWordAddress, 1,
		d.StatusWordAddress)
}

func (d *DoorUnlockControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]电磁门解锁,发送(%s,%d),状态字地址%s",
		d.ControlWordAddress, 1,
		d.StatusWordAddress)
}

// TemperatureControlAction 温度控制动作
type TemperatureControlAction struct {
	*ControlAction
	TargetTemperature *data.AddressValue // 目标温度，单位：0.1摄氏度
}

func NewTemperatureControlAction(controlWordAddress string, statusWordAddress string,
	targetTemperature *data.AddressValue) *TemperatureControlAction {
	temperatureControlAction := &TemperatureControlAction{
		ControlAction:     NewControlAction(controlWordAddress, statusWordAddress),
		TargetTemperature: targetTemperature,
	}
	temperatureControlAction.ControlAction.AddressValueList = append(temperatureControlAction.ControlAction.AddressValueList,
		targetTemperature,
		data.NewAddressValue(data.TEMPERATURE_UPPER_VALUE_ADDRESS, 2200),
		data.NewAddressValue(data.TEMPERATURE_LOWER_VALUE_ADDRESS, 0))
	return temperatureControlAction
}

func (c *TemperatureControlAction) BeforeExecuteInfo() string {
	return fmt.Sprintf("[开始]温度控制,目标温度%.1f℃,发送(%s,%d),(%s,%d),状态字地址%s",
		float32(c.TargetTemperature.Value)/10,
		c.ControlWordAddress, 1,
		c.TargetTemperature.Address, c.TargetTemperature.Value,
		c.StatusWordAddress)
}

func (c *TemperatureControlAction) AfterExecuteInfo() string {
	return fmt.Sprintf("[结束]温度控制,目标温度%1.f℃,发送(%s,%d),(%s,%d),状态字地址%s",
		float32(c.TargetTemperature.Value)/10,
		c.ControlWordAddress, 1,
		c.TargetTemperature.Address, c.TargetTemperature.Value,
		c.StatusWordAddress)
}
