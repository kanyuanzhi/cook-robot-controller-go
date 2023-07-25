package grpc

import (
	"context"
	"cook-robot-controller-go/core"
	"cook-robot-controller-go/data"
	pb "cook-robot-controller-go/grpc/commandRPC"
	"cook-robot-controller-go/instruction"
	"cook-robot-controller-go/logger"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
)

type command struct {
	pb.UnimplementedCommandServiceServer
	controller *core.Controller
}

func (c *command) Execute(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	var commandMap map[string]interface{}
	err := json.Unmarshal([]byte(req.GetCommandJson()), &commandMap)
	if err != nil {
		return nil, fmt.Errorf("无法解析命令JSON：%v", err)
	}
	if commandMap["commandType"].(string) == data.COMMAND_TYPE_MULTIPLE {
		// 多指令型命令，需要判断是否有其他命令在运行
		if c.controller.CurrentCommandName == "" {
			// 无命令在运行
			c.controller.CurrentCommandName = commandMap["commandName"].(string)
			if commandMap["commandName"].(string) == "cook" {
				c.controller.IsCooking = true
				c.controller.CurrentDishUuid = commandMap["dishUuid"].(string)
			}
			for _, ins := range commandMap["instructions"].([]interface{}) {
				//logger.Log.Println(ins.(map[string]interface{}))
				var instructionType = instruction.InstructionType((ins.(map[string]interface{})["instructionType"]).(string))
				instructionStruct := instruction.InstructionTypeToStruct[instructionType]
				err := mapstructure.Decode(ins, &instructionStruct)
				if err != nil {
					logger.Log.Println(err)
				}
				instructionStruct.AddToController(c.controller)
			}
			go c.controller.Start()
		} else {
			// 有命令在运行
			logger.Log.Printf("%s命令正在运行，无法执行当前命令", c.controller.CurrentCommandName)
			return &pb.CommandResponse{Result: 0}, nil
		}
	} else {
		// 单指令，立即运行
		instructionInter := commandMap["instructions"].([]interface{})[0]
		var instructionType = instruction.InstructionType((instructionInter.(map[string]interface{})["instructionType"]).(string))
		instructionStruct := instruction.InstructionTypeToStruct[instructionType]
		err := mapstructure.Decode(instructionInter, &instructionStruct)
		if err != nil {
			logger.Log.Println(err)
		}
		instructionStruct.ExecuteImmediately(c.controller)
	}

	return &pb.CommandResponse{Result: 1}, nil
}

func (c *command) FetchStatus(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	type ControllerStatus struct {
		CurrentCommandName              string                `json:"currentCommandName"`
		CurrentDishUuid                 string                `json:"currentDishUuid"`
		CurrentInstructionName          string                `json:"currentInstructionName"`
		CurrentInstructionInfo          *data.InstructionInfo `json:"currentInstructionInfo"`
		IsPausing                       bool                  `json:"isPausing"`
		IsRunning                       bool                  `json:"isRunning"`
		IsCooking                       bool                  `json:"isCooking"`
		IsPausingWithMovingFinished     bool                  `json:"isPausingWithMovingFinished"`
		IsPausingWithMovingBackFinished bool                  `json:"isPausingWithMovingBackFinished"`
		IsStirFrying                    bool                  `json:"isStirFrying"`
		BottomTemperature               uint32                `json:"bottomTemperature"`
		InfraredTemperature             uint32                `json:"infraredTemperature"`
		CookingTime                     int64                 `json:"cookingTime"`
	}
	controllerStatus := ControllerStatus{
		CurrentCommandName:              c.controller.CurrentCommandName,
		CurrentDishUuid:                 c.controller.CurrentDishUuid,
		CurrentInstructionInfo:          c.controller.CurrentInstructionInfo,
		IsPausing:                       c.controller.IsPausing,
		IsRunning:                       c.controller.IsRunning,
		IsCooking:                       c.controller.IsCooking,
		IsPausingWithMovingFinished:     c.controller.IsPausingWithMovingFinished,
		IsPausingWithMovingBackFinished: c.controller.IsPausingWithMovingBackFinished,
		IsStirFrying:                    c.controller.IsPausePermitted,
		BottomTemperature:               c.controller.TcpServer.RealtimeValueMap[data.TEMPERATURE_BOTTOM_ADDRESS],
		InfraredTemperature:             c.controller.TcpServer.RealtimeValueMap[data.TEMPERATURE_INFRARED_ADDRESS],
		CookingTime:                     c.controller.CookingTime,
	}
	statusJSON, _ := json.Marshal(controllerStatus)
	return &pb.FetchResponse{StatusJson: string(statusJSON)}, nil
}

func (c *command) Pause(ctx context.Context, req *pb.PauseAndResumeRequest) (*pb.PauseAndResumeResponse, error) {
	c.controller.Pause()
	return &pb.PauseAndResumeResponse{Result: 1}, nil
}

func (c *command) Resume(ctx context.Context, req *pb.PauseAndResumeRequest) (*pb.PauseAndResumeResponse, error) {
	c.controller.Resume()
	return &pb.PauseAndResumeResponse{Result: 1}, nil
}

func (c *command) Shutdown(ctx context.Context, req *pb.ShutdownRequest) (*pb.ShutdownResponse, error) {
	os.Exit(1)
	return &pb.ShutdownResponse{Result: 1}, nil
}
