package grpc

import (
	"context"
	"cook-robot-controller-go/core"
	pb "cook-robot-controller-go/grpc/commandRPC"
	"cook-robot-controller-go/instruction"
	"cook-robot-controller-go/logger"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type command struct {
	pb.UnimplementedCommandServiceServer
	controller *core.Controller
}

func (c *command) Execute(ctx context.Context, req *pb.MultipleCommandRequest) (*pb.CommandResponse, error) {
	var commandMap map[string]interface{}
	err := json.Unmarshal([]byte(req.GetCommandJson()), &commandMap)
	if err != nil {
		return nil, fmt.Errorf("无法解析命令JSON：%v", err)
	}
	for _, ins := range commandMap["instructions"].([]interface{}) {
		//logger.Log.Println(ins.(map[string]interface{}))
		var instructionType = instruction.InstructionType((ins.(map[string]interface{})["instructionType"]).(string))
		instructionStruct := instruction.InstructionTypeToStruct[instructionType]
		err := mapstructure.Decode(ins, &instructionStruct)
		if err != nil {
			logger.Log.Println(err)
		}
		//logger.Log.Println(instructionStruct)
		instructionStruct.AddToController(c.controller)
	}
	go c.controller.Start()

	return &pb.CommandResponse{Result: 1}, nil
}

func (c *command) Run(ctx context.Context, req *pb.SingleCommandRequest) (*pb.CommandResponse, error) {
	go c.controller.Start()
	return &pb.CommandResponse{Result: 1}, nil
}

func (c *command) Fill(ctx context.Context, req *pb.MultipleCommandRequest) (*pb.CommandResponse, error) {
	var commandMap map[string]interface{}
	err := json.Unmarshal([]byte(req.GetCommandJson()), &commandMap)
	if err != nil {
		return nil, fmt.Errorf("无法解析命令JSON：%v", err)
	}
	for _, ins := range commandMap["instructions"].([]interface{}) {
		//logger.Log.Println(ins.(map[string]interface{}))
		var instructionType = instruction.InstructionType((ins.(map[string]interface{})["instructionType"]).(string))
		instructionStruct := instruction.InstructionTypeToStruct[instructionType]
		err := mapstructure.Decode(ins, &instructionStruct)
		if err != nil {
			logger.Log.Println(err)
		}
		//logger.Log.Println(instructionStruct)
		instructionStruct.AddToController(c.controller)
	}
	return &pb.CommandResponse{Result: 1}, nil
}
