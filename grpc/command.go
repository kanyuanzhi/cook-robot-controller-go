package grpc

import (
	"context"
	pb "cook-robot-controller-go/grpc/commandServer"
	"cook-robot-controller-go/instruction"
	"cook-robot-controller-go/logger"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type commandServer struct {
	pb.UnimplementedCommandServiceServer
}

func (c *commandServer) Execute(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	//var command command.Command
	var com map[string]interface{}
	err := json.Unmarshal([]byte(req.GetCommandJson()), &com)
	if err != nil {
		return nil, fmt.Errorf("无法解析命令JSON：%v", err)
	}
	logger.Log.Println(com)
	logger.Log.Println(com["instructions"])
	for _, ins := range com["instructions"].([]interface{}) {
		logger.Log.Println(ins)
		logger.Log.Println(ins.(map[string]interface{})["instruction_type"])
		if ins.(map[string]interface{})["instruction_type"] == "alex" {
			var ins1 instruction.AlexInstruction
			mapstructure.Decode(ins, &ins1)
			logger.Log.Println(ins1)
		}
	}

	return &pb.CommandResponse{Result: 1}, nil
}
