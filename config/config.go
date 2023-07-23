package config

import (
	"cook-robot-controller-go/logger"
	"github.com/spf13/viper"
)

var App *AppConfig

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

func init() {
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
}
