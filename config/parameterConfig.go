package config

import (
	"cook-robot-controller-go/utils"
)

var Parameter = &ParameterConfig{}

type XAxisConfig struct {
	MoveSpeed          uint32 `mapstructure:"moveSpeed"`
	WithdrawerPosition uint32 `mapstructure:"withdrawerPosition"`
	ReadyPosition      uint32 `mapstructure:"readyPosition"`
	Box1Position       uint32 `mapstructure:"box1Position"`
	Box2Position       uint32 `mapstructure:"box2Position"`
	Box3Position       uint32 `mapstructure:"box3Position"`
	Box4Position       uint32 `mapstructure:"box4Position"`
	SafePosition       uint32 `mapstructure:"safePosition"`
	ShakeAmount        uint32 `mapstructure:"shakeAmount"`
	ShakeUpwardSpeed   uint32 `mapstructure:"shakeUpwardSpeed"`
	ShakeDownwardSpeed uint32 `mapstructure:"shakeDownwardSpeed"`
	ShakeDistance      uint32 `mapstructure:"shakeDistance"`
}

type YAxisConfig struct {
	MoveSpeed               uint32 `mapstructure:"moveSpeed"`
	DishOutAmount           uint32 `mapstructure:"dishOutAmount"`
	DishOutUpwardSpeed      uint32 `mapstructure:"dishOutUpwardSpeed"`
	DishOutDownwardSpeed    uint32 `mapstructure:"dishOutDownwardSpeed"`
	DishOutHighPosition     uint32 `mapstructure:"dishOutHighPosition"`
	DishOutLowPosition      uint32 `mapstructure:"dishOutLowPosition"`
	IngredientPosition      uint32 `mapstructure:"ingredientPosition"`
	LiquidSeasoningPosition uint32 `mapstructure:"liquidSeasoningPosition"`
	SolidSeasoningPosition  uint32 `mapstructure:"solidSeasoningPosition"`
	StirFry1Position        uint32 `mapstructure:"stirFry1Position"`
	StirFry2Position        uint32 `mapstructure:"stirFry2Position"`
	StirFry3Position        uint32 `mapstructure:"stirFry3Position"`
	WashPosition            uint32 `mapstructure:"washPosition"`
	PourPosition            uint32 `mapstructure:"pourPosition"`
}

type R1AxisConfig struct {
	MaxRotateSpeed uint32 `mapstructure:"maxRotateSpeed"`
}

type WashConfig struct {
	RotateSpeed       uint32 `mapstructure:"rotateSpeed"`
	RotateCrossAmount uint32 `mapstructure:"rotateCrossAmount"`
	AddWaterDuration  uint32 `mapstructure:"addWaterDuration"`
	Temperature       uint32 `mapstructure:"temperature"`
	Duration          uint32 `mapstructure:"duration"`
	PumpNumber        uint32 `mapstructure:"pumpNumber"`
	PourWaterDuration uint32 `mapstructure:"pourWaterDuration"`
}

type LampblackPurify struct {
	Enable           bool   `mapstructure:"enable"`
	VentingMode      uint32 `mapstructure:"ventingMode"`
	PurificationMode uint32 `mapstructure:"purificationMode"`
	AutoStart        bool   `mapstructure:"autoStart"`
	AutoStartMode    uint32 `mapstructure:"autoStartMode"`
}

type ParameterConfig struct {
	XAxis           XAxisConfig     `mapstructure:"xAxis"`
	YAxis           YAxisConfig     `mapstructure:"yAxis"`
	R1Axis          R1AxisConfig    `mapstructure:"r1Axis"`
	Wash            WashConfig      `mapstructure:"wash"`
	LampblackPurify LampblackPurify `mapstructure:"lampblackPurify"`
}

func (w *ParameterConfig) Reload() {
	utils.Reload("parameterConfig", Parameter)
}
