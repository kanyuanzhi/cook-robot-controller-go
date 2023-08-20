package config

import (
	"cook-robot-controller-go/logger"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

var App *AppConfig
var Parameter *ParameterConfig

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type ModbusConfig struct {
	TargetHost string `mapstructure:"targetHost"`
	TargetPort uint16 `mapstructure:"targetPort"`
}

type AppConfig struct {
	DebugMode bool
	GRPC      GRPCConfig   `mapstructure:"grpc"`
	Modbus    ModbusConfig `mapstructure:"modbus"`
}

func (m *AppConfig) Reload() {
	viper.SetConfigName("controllerConfig")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Println("无法读取配置文件:", err)
		return
	}

	err = viper.Unmarshal(App)
	if err != nil {
		logger.Log.Println("解析配置文件失败:", err)
		return
	}
}

type XAxisConfig struct {
	MoveSpeed          uint32 `mapstructure:"moveSpeed"`
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
	Enable bool   `mapstructure:"enable"`
	Mode   uint32 `mapstructure:"mode"`
}

type ParameterConfig struct {
	XAxis           XAxisConfig     `mapstructure:"xAxis"`
	YAxis           YAxisConfig     `mapstructure:"yAxis"`
	Wash            WashConfig      `mapstructure:"wash"`
	LampblackPurify LampblackPurify `mapstructure:"lampblackPurify"`
}

func (w *ParameterConfig) Reload() {
	viper.SetConfigName("parameterConfig")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Println("无法读取配置文件:", err)
		return
	}

	err = viper.Unmarshal(Parameter)
	if err != nil {
		logger.Log.Println("解析配置文件失败:", err)
		return
	}
}

func init() {
	file, err := os.OpenFile("controller.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Unable to create log file:", err)
	}
	//defer file.Close()

	// 使用io.MultiWriter将日志输出同时写入控制台和文件
	logWriter := io.MultiWriter(os.Stdout, file)

	logger.Log = log.New(logWriter, "", log.Lmicroseconds)

	App = &AppConfig{
		DebugMode: false,
		GRPC: GRPCConfig{
			Host: "localhost",
			Port: 50051,
		},
		Modbus: ModbusConfig{
			TargetHost: "192.168.6.6",
			TargetPort: 502,
		},
	}
	App.Reload()

	Parameter = &ParameterConfig{
		XAxis: XAxisConfig{
			MoveSpeed:          20000,
			ReadyPosition:      54396,
			Box1Position:       42690,
			Box2Position:       25937,
			Box3Position:       9069,
			Box4Position:       20,
			SafePosition:       10,
			ShakeAmount:        5,
			ShakeUpwardSpeed:   30000,
			ShakeDownwardSpeed: 20000,
			ShakeDistance:      3200,
		},
		YAxis: YAxisConfig{
			MoveSpeed:               600,
			DishOutAmount:           3,
			DishOutUpwardSpeed:      600,
			DishOutDownwardSpeed:    2200,
			DishOutHighPosition:     4200,
			DishOutLowPosition:      4700,
			IngredientPosition:      2000,
			LiquidSeasoningPosition: 3200,
			SolidSeasoningPosition:  3485,
			StirFry1Position:        3500,
			StirFry2Position:        3500,
			StirFry3Position:        3500,
			WashPosition:            3400,
			PourPosition:            10,
		},
		Wash: WashConfig{
			RotateSpeed:       1500,
			RotateCrossAmount: 3,
			AddWaterDuration:  15000,
			Temperature:       900,
			Duration:          60000,
			PumpNumber:        7,
			PourWaterDuration: 4000,
		},
		LampblackPurify: LampblackPurify{
			Enable: true,
			Mode:   2,
		},
	}
	Parameter.Reload()
}
