package config

import (
	"cook-robot-controller-go/utils"
)

var App = &AppConfig{}

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
	utils.Reload("controllerConfig", App)
}
