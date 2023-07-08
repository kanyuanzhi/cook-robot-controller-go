package grpc

import (
	"context"
	"cook-robot-controller-go/core"
	pb "cook-robot-controller-go/grpc/commandServer"
	"cook-robot-controller-go/instruction"
	"cook-robot-controller-go/logger"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type commandServer struct {
	pb.UnimplementedCommandServiceServer
	controller *core.Controller
}

func (c *commandServer) Execute(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	var commandMap map[string]interface{}
	err := json.Unmarshal([]byte(req.GetCommandJson()), &commandMap)
	if err != nil {
		return nil, fmt.Errorf("无法解析命令JSON：%v", err)
	}
	for _, ins := range commandMap["instructions"].([]interface{}) {
		//logger.Log.Println(ins.(map[string]interface{})["instruction_type"])
		var instructionType = instruction.InstructionType((ins.(map[string]interface{})["instruction_type"]).(string))
		instructionStruct := instruction.InstructionTypeToStruct[instructionType]
		err := mapstructure.Decode(ins, &instructionStruct)
		if err != nil {
			logger.Log.Println(err)
		}
		//logger.Log.Println(instructionStruct)
		instructionStruct.AddToController(c.controller)
	}
	go c.controller.Run()

	return &pb.CommandResponse{Result: 1}, nil
}
